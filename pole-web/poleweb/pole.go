package poleweb

import (
	"log"
	"net/http"
)

// HandlerFunc方法实现
// type HandlerFunc func(http.ResponseWriter, *http.Request)
// 参数修改为context类型
type HandlerFunc func(*Context)

// 实现ServeHTTP方法即实现了接口
type Engine struct {
	*RouterGroup
	groups []*RouterGroup // store all groups
	router *router
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 中间件支持
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// 结构体的构造器
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 添加方法到路由表中
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// 定义GET请求的方法并添加到engine的HandlerFunc方法字典中
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// 定义POST请求的方法并添加到engine的HandlerFunc方法字典中
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 开启一个http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
