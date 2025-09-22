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

func TestGetUser(t *testing.T) {
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

func loginUser(
	t *testing.T,
	server *httptest.Server,
	username string,
	passwword string,
) string {
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

	return jwtToken.Token
}

func TestCreateUser(t *testing.T) {
	mux := Setup(t)
	server := httptest.NewServer(mux)
	defer server.Close()

	token := loginUser(t, server, "admin", "admin")

	newUser := schemas.CreateUserDTO{
		UserName:  "tyler.durden",
		FirstName: "tyler",
		LastName:  "durden",
		Email:     "tyler@paperstreetsoap.com",
		Password:  "log-me-in",
	}
	payload, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("could not marshal create user data: %v", err)
	}

	req, err := http.NewRequest("POST", server.URL+"/", bytes.NewBuffer(payload))

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "appliction/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, but got: %d", resp.StatusCode)
	}
}
