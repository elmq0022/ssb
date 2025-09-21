//go:build integration
// +build integration

package integration

import (
	"testing"
	"net/http/httptest"
	"net/http"
)

func TestHealthzIntegration(t *testing.T) {
	mux := Setup(t)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
}
