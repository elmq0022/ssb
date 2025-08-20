package healthz

import (
	"net/http"
	"ssb/internal/pkg/router"
)

func NewRouter() *router.Router {
	r := router.NewRouter()
	r.Get("/", healthzHandler)
	return r
}

func healthzHandler(r *http.Request) (any, int, error) {
	return map[string]string{"status": "ok"}, http.StatusOK, nil
}
