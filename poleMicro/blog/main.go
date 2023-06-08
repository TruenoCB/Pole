package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TruenoCB/poleweb"
)

type User struct {
	Name      string   `xml:"name" json:"name" polevalidate:"required"`
	Age       int      `xml:"age" json:"age" validate:"required,max=50,min=18"`
	Addresses []string `json:"addresses"`
	Email     string   `json:"email" polevalidate:"required"`
}

func Log(next poleweb.HandlerFunc) poleweb.HandlerFunc {
	return func(ctx *poleweb.Context) {
		log.Println("打印请求参数")
		next(ctx)
		log.Println("返回执行时间")
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
		user := User{Name: "jack"}
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

	test.Post("/files", func(ctx *poleweb.Context) {
		m, _ := ctx.GetPostFormMap("user")
		files := ctx.FormFiles("file")
		for _, file := range files {
			ctx.SaveUploadedFile(file, "blog/upload/"+file.Filename)
		}
		ctx.JSON(http.StatusOK, m)
	})

	test.Post("/jsonParam", func(ctx *poleweb.Context) {
		user := make([]User, 0)
		ctx.DisallowUnknownFields = true
		//ctx.IsValidate = true
		err := ctx.BindJson(&user)
		if err == nil {
			ctx.JSON(http.StatusOK, user)
		} else {
			log.Println(err)
		}
	})

	engine.Run()
}
