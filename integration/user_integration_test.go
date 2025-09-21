//go:build integration
// +build integration

package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserIntegration(t *testing.T) {
	mux := Setup(t)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/users/admin")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d, but %d", http.StatusOK, resp.StatusCode)
	}
}
