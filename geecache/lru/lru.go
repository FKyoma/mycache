package lru

import "container/list"

// Cache是LRU缓存，目前并发访问不安全！
type Cache struct {
	maxBytes int64      //允许使用的最大的内存
	nbytes   int64      //nowbytes 目前已使用的内存
	ll       *list.List //双向链表
	cache    map[string]*list.Element
	//键是字符串，值是双向链表中对应节点的指针

	//OnEvicted 是某条记录被移除时的回调函数，可以为nil
	OnEvicted func(key string, value Value)
}

//entry是双向链表的数据结构
type entry struct {
	key   string
	value Value
}

//用于返回值所占用的内存大小
type Value interface {
	Len() int
}

//初始化
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//增加
func (c *Cache) Add(key string, value Value) {

	//首先查找是否重复，如果重复则更新值
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry) //临时存储值

		c.nbytes += int64(value.Len()) - int64(kv.value.Len()) //数据更新后，nbytes改变
		kv.value = value
	} else {
		//未重复则添加新键值对
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		//如果当前内存超过了设定值时，淘汰最近最少访问节点
		c.RemoveOldest()
	}
}

//查找 第一步从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		//约定Front为队尾,Back为队首
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//删除
func (c *Cache) RemoveOldest() {

	ele := c.ll.Back()
	if ele != nil {
		//删除队首
		c.ll.Remove(ele)

		kv := ele.Value.(*entry)

		delete(c.cache, kv.key)

		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
