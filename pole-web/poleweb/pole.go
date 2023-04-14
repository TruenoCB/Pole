package poleweb

import (
	"fmt"
	"net/http"
)

// HandlerFunc方法实现
type HandlerFunc func(http.ResponseWriter, *http.Request)

// 实现ServeHTTP方法即实现了接口
type Engine struct {
	router map[string]HandlerFunc
}

// 结构体的构造器
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加方法到路由表中
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// 定义GET请求的方法并添加到engine的HandlerFunc方法字典中
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// 定义POST请求的方法并添加到engine的HandlerFunc方法字典中
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 开启一个http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok { //根据请求方法和路径生成key，到路由表中寻找对应的HandlerFunc方法
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
