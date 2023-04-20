package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"poleweb"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func onlyForV3() poleweb.HandlerFunc {
	return func(c *poleweb.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v3", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := poleweb.New()
	r.Use(poleweb.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	r.GET("/", func(c *poleweb.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/hello", func(c *poleweb.Context) {
		// url后带上参数
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *poleweb.Context) {
		c.JSON(http.StatusOK, poleweb.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/hello/:name", func(c *poleweb.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *poleweb.Context) {
			c.HTML(http.StatusOK, "css.tmpl", nil)
		})

		v1.GET("/hello", func(c *poleweb.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *poleweb.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *poleweb.Context) {
			c.JSON(http.StatusOK, poleweb.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}
	v3 := r.Group("/v3")
	v3.Use(onlyForV3()) // v3 group middleware
	{
		v3.GET("/hello/:name", func(c *poleweb.Context) {
			// expect /hello/trueno
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	stu1 := &student{Name: "trueno", Age: 23}
	stu2 := &student{Name: "Jack", Age: 22}

	r.GET("/students", func(c *poleweb.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", poleweb.H{
			"title":  "pole",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *poleweb.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", poleweb.H{
			"title": "pole",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":1018")
}
