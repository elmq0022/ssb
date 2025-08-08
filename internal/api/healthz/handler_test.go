package api_test

import (
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/healthz"
	"testing"
)

func TestHealthzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	r := api.R
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	expectedBody := `{"status":"ok"}` + "\n"
	if w.Body.String() != expectedBody {
		t.Errorf("unexpected body: got %q, want %q", w.Body.String(), expectedBody)
	}
}
