package router

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONHandler func(r *http.Request) (any, int, error)

func jsonToHttpHandler(h JSONHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val, status, err := h(r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(val)
	}
}

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) Handle(method, path string, jsonHandler JSONHandler) {
	route := fmt.Sprintf("%s %s", method, path)
	httpHandler := jsonToHttpHandler(jsonHandler)
	r.mux.HandleFunc(route, httpHandler)
}

func (r *Router) Get(path string, handler JSONHandler) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler JSONHandler) {
	r.Handle(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler JSONHandler) {
	r.Handle(http.MethodPut, path, handler)
}

func (r *Router) Delete(path string, handler JSONHandler) {
	r.Handle(http.MethodDelete, path, handler)
}

func (r *Router) Serve(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
