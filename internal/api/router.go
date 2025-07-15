package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", healthzHandler)
	return r
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aaplication/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
