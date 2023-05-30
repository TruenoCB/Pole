package poleweb

import (
	"fmt"
	"net/http"
)

type engine struct {
}

func New() *engine {
	return &engine{}
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n", r.Method)
	w.Write([]byte(r.URL.Path))
}

func (e *engine) Run() {
	
	http.Handle("/", e)
	http.HandleFunc("/abc", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("helloworld"))
	})
	http.ListenAndServe(":8080", nil)
}
