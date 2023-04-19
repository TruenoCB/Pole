package main

import (
	"log"
	"net/http"
	"poleweb"
	"time"
)

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
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *poleweb.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
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

	r.Run(":1018")
}
