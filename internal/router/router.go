package router

import (
	"fmt"
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	route := fmt.Sprintf("%s %s", method, path)
	r.mux.HandleFunc(route, handler)
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}

func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}

func (r *Router) Serve(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
