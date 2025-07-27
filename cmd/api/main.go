package main

import (
	"log"
	"net/http"
	"ssb/internal/api/healthz"
)

func main() {
	r := api.NewRouter()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
