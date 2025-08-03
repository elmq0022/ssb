package router_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"ssb/internal/router"
	"testing"
)

func TestRouterGet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	}
	req := httptest.NewRequest(http.MethodGet, "/home", nil)
	w := httptest.NewRecorder()

	r := router.NewRouter()
	r.Get("/home", handler)
	r.Serve(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("%v", r.Routes)
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	if w.Body.String() != "hello" {
		t.Fatalf("expected response body 'hello', but got %q", w.Body.String())
	}
}
