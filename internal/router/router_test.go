package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ssb/internal/router"
	"testing"
)

func TestRouterGet(t *testing.T) {
	result := map[string]string{"msg": "hello"}
	handler := func(r *http.Request) (any, int, error) {
		return result, 200, nil
	}
	r := router.NewRouter()
	r.Get("/home", handler)

	req := httptest.NewRequest(http.MethodGet, "/home", nil)
	w := httptest.NewRecorder()
	r.Serve(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
	var buff bytes.Buffer
	json.NewEncoder(&buff).Encode(result)
	want := buff.String()
	got := w.Body.String()
	if want != got {
		t.Fatalf("expected response body %q, but got %q", want, w.Body.String())
	}
}

func TestRouterPost(t *testing.T) {
	t.Fatalf("Fails")
}

func TestRouterPut(t *testing.T) {
	t.Fatalf("Fails")
}

func TestRouterDelete(t *testing.T) {
	t.Fatalf("Fails")
}
