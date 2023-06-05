package poleweb

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/TruenoCB/poleweb/internal/utils"
	"github.com/TruenoCB/poleweb/render"
)

type HandlerFunc func(ctx *Context)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

type Engine struct {
	router
	funcMap     template.FuncMap
	HTMLRender  render.HTMLRender
	middles     []MiddlewareFunc
	pool        sync.Pool
	OpenGateway bool
}

func New() *Engine {
	engine := &Engine{
		router: router{},
	}
	engine.router.engine = engine
	engine.pool.New = func() any {
		return engine.allocateContext()
	}
	return engine
}

func (e *Engine) allocateContext() any {
	return &Context{engine: e}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.pool.Get().(*Context)
	ctx.W = w
	ctx.R = r
	e.httpRequestHandle(ctx, w, r)
	e.pool.Put(ctx)
}

func (e *Engine) httpRequestHandle(ctx *Context, w http.ResponseWriter, r *http.Request) {
	method := r.Method
	for _, group := range e.routerGroups {
		routerName := utils.SubStringLast(r.URL.Path, "/"+group.name)
		// get/1
		node := group.treeNode.Get(routerName)
		if node != nil && node.isEnd {
			/* ctx := &Context{   已被sync.Pool池化
				W:      w,
				R:      r,
				engine: e,
			} */
			//路由匹配上了
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				group.methodHandle(node.routerName, ANY, handle, ctx)
				return
			}
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				group.methodHandle(node.routerName, method, handle, ctx)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "%s %s not allowed \n", r.RequestURI, method)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s  not found \n", r.RequestURI)
}

func (e *Engine) Run() {
	http.Handle("/", e)
	err := http.ListenAndServe(":8600", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadTemplate(pattern string) {
	t := template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
	e.SetHtmlTemplate(t)
}
func (e *Engine) SetHtmlTemplate(t *template.Template) {
	e.HTMLRender = render.HTMLRender{Template: t}
}
