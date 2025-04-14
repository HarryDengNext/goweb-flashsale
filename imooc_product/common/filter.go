package common

import (
	"net/http"
	"strings"
)

// 声明一个新的数据类型（函数类型）
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// 拦截器结构体
type Filter struct {
	filterMap map[string]FilterHandle
}

// Filter初始化函数
func NewFilter()*Filter  {
	return &Filter{filterMap:make(map[string]FilterHandle)}
}

func (f *Filter)RegisterFilterUrl(url string, handler FilterHandle)  {
	f.filterMap[url] = handler
}

// 根据url获取对应的handle
func (f *Filter) GetFilterHandle(url string) FilterHandle  {
	return f.filterMap[url]
}

type WebHandle func(rw http.ResponseWriter, req *http.Request)

// 执行拦截器， 返回函数类型
func (f *Filter) Handle (webhandle WebHandle) func(rw http.ResponseWriter, r *http.Request)  {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if strings.Contains(r.RequestURI, path){
				// 执行拦截业务逻辑
				err := handle(rw, r)
				if err!=nil {
					rw.Write([]byte(err.Error()))
					return
				}
				// 跳出循环
				break
			}
		}
		// 执行正常注册的函数
		webhandle(rw, r)
	}	
}

