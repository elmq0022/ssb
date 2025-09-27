//go:build integration
// +build integration

package testutil

// TODO: move utility functions here.
import (
	"testing"
	"encoding/json"
	"net/http/httptest"
	"net/http"
	"ssb/internal/schemas"
	"bytes"
	"io"
)

func LoginUser(
	t *testing.T,
	server *httptest.Server,
	username string,
	password string,
) string {
	data := schemas.LoginRequest{
		Username: username,
		Password: password,
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

	return jwtToken.Token
}


func MakeRequest(
	t *testing.T, token, method, url string,
	payload io.Reader) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		t.Fatalf("%v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	return req
}
