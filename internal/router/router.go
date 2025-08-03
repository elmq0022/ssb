package router

import (
	"fmt"
	"net/http"
	"sort"
)

type Router struct {
	// HttpMethod -> Path -> Handler Function
	Routes map[string]map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		Routes: make(map[string]map[string]http.HandlerFunc),
	}
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	if r.Routes[method] == nil {
		r.Routes[method] = make(map[string]http.HandlerFunc)
	}
	r.Routes[method][path] = handler
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
	if methodRoutes, ok := r.Routes[req.Method]; ok {
		if handler, ok := methodRoutes[req.URL.Path]; ok {
			handler(w, req)
			return
		}
	}

	http.NotFound(w, req)
}

func (r *Router) PrintRoutes() []string{
	var results []string
	for method, routes := range r.Routes {
		for path := range routes {
			results = append(results, fmt.Sprintf("%s %s", method, path))	
		}
	}
	sort.Strings(results)
	return results
}

