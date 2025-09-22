//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ssb/internal/schemas"
	"testing"
)

func TestAuthIntegration(t *testing.T) {
	mux := Setup(t)
	server := httptest.NewServer(mux)
	defer server.Close()

	data := schemas.LoginRequest{
		Username: "admin",
		Password: "admin",
	}

	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(
		server.URL+"/auth/login",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		t.Fatalf("got error from POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wanted 200 OK, but got %d", resp.StatusCode)
	}

	var jwtToken schemas.JsonToken
	if err := json.NewDecoder(resp.Body).Decode(&jwtToken); err != nil {
		t.Fatalf("failed to decode jwtToken: %v", err)
	}
}
