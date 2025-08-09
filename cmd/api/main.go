package main

import (
	"log"
	"net/http"
	"ssb/internal/api/articles"
	"ssb/internal/api/healthz"
	"ssb/internal/router"
)

func main() {
	mux := router.NewRouter()
	mux.Mount("/healthz", healthz.NewRouter())
	mux.Mount("/articles", articles.NewRouter())
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
