//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"
)

func TestHealthzEnpoint(t *testing.T) {
	server, _, _ := Setup(t)

	resp, err := http.Get(server.URL + "/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
}
