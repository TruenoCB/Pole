package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TruenoCB/poleweb"
)

type User struct {
	Name string
}

func Log(next poleweb.HandlerFunc) poleweb.HandlerFunc {
	return func(ctx *poleweb.Context) {
		fmt.Println("打印请求参数")
		next(ctx)
		fmt.Println("返回执行时间")
	}
}

func main() {
	engine := poleweb.New()
	engine.LoadTemplate("blog/tpl/*.html")
	test := engine.Group("test")
	test.Any("/t1", func(ctx *poleweb.Context) {
		log.Println("handler")
		fmt.Fprintf(ctx.W, "%s hello world", ctx.R.RemoteAddr)
	})
	test.Any("/t1/:id/car", func(ctx *poleweb.Context) {
		log.Println("文件请求")
		ctx.FileFromFS("car.jpg", http.Dir("blog/tpl"))
	})
	test.Any("/login", func(ctx *poleweb.Context) {
		user := User{"jack"}
		log.Println("render模板接口")
		err := ctx.Template("login.html", user)
		if err != nil {
			log.Println(err)
		}
	})
	test.Any("/index", func(ctx *poleweb.Context) {
		log.Println("render模板接口")
		ctx.HTML(http.StatusOK, "login.html")
	})
	test.Any("/mid", func(ctx *poleweb.Context) {
		log.Println("中间件链式调用")
		ctx.HTML(http.StatusOK, "<h1>中间件链式调用</h1>")
	}, Log)
	engine.Run()
}
