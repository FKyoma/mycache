package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 规定访问路径形式为  /<basepath>/<groupname>/<key>
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	//parts 字符串切片 分别为  parts[1]= <groupname>  parts[2]= <key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)

	//如果返回的parts没有两个,那么请求格式错误

	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	//parts[0]和parts[1]分别代表groupname和key
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//定义网络文件的类型和网页的编码        "application/octet-stream" 为二进制流数据(常见的文件下载)
	w.Header().Set("Content-Type", "application/octet-stream")
	//写入值
	w.Write(view.ByteSlice())
}
