package api

import (
	"net/http"
	"ssb/internal/router"
)

var R *router.Router = router.NewRouter()

func init() {
	R.Get("/healthz", healthzHandler)
}

func healthzHandler(r *http.Request) (any, int, error) {
	return map[string]string{"status": "ok"}, http.StatusOK, nil
}
