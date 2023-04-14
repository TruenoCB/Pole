package main

import (
	"fmt"
	"net/http"
	"poleweb"
)

func main() {
	r := poleweb.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		//打印请求头部信息
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.GET("/pole", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "its a concise frame")

	})

	r.Run(":1018")
}
