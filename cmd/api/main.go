package main

import (
	"log"
	"net/http"
	"ssb/internal/api/healthz"
	"ssb/internal/router"
)

func main() {
	mux := router.NewRouter()
	mux.Mount("/healthz", healthz.NewRouter())
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
