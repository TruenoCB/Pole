package poleweb

import (
	"log"
	"net/http"
)

// HandlerFunc方法实现
//type HandlerFunc func(http.ResponseWriter, *http.Request)
//参数修改为context类型
type HandlerFunc func(*Context)

// 实现ServeHTTP方法即实现了接口
type Engine struct {
	router *router
}

// 结构体的构造器
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 添加方法到路由表中
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
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
	c := newContext(w, req)
	engine.router.handle(c)
}
