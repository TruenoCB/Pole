package main

import (
	"net/http"
	"poleweb"
)

func main() {
	r := poleweb.New()
	r.GET("/", func(c *poleweb.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Pole</h1>")
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

	r.GET("/assets/*filepath", func(c *poleweb.Context) {
		c.JSON(http.StatusOK, poleweb.H{"filepath": c.Param("filepath")})
	})

	r.Run(":1018")
}
