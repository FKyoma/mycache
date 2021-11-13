package geecache


type ByteView struct {
	b []byte
}

//返回长度
func (v ByteView) Len() int {
	return len(v.b)
}

//返回一个拷贝，防止缓存被修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

//String以字符串形式返回数据，必要时制作副本。
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
