//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"ssb/internal/schemas"
	"testing"
)

func TestGetUser(t *testing.T) {
	server := Setup(t)

	resp, err := http.Get(server.URL + "/users/admin")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d, but %d", http.StatusOK, resp.StatusCode)
	}
}

func makeRequest(
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

func loginUser(
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

func TestCreateUser(t *testing.T) {
	server := Setup(t)

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

	req, err := http.NewRequest(
		"POST",
		server.URL+"/users/",
		bytes.NewBuffer(payload),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("want 201, but got: %d", resp.StatusCode)
	}
}

func TestUpdateUser(t *testing.T) {
	server := Setup(t)

	token := loginUser(t, server, "admin", "admin")

	updatedEmail := "narrator@paperstreetsoap.com"
	updateUserData := schemas.UpdateUserDTO{
		Email: &updatedEmail,
	}
	payload, err := json.Marshal(updateUserData)
	if err != nil {
		t.Fatalf("%v", err)
	}

	req := makeRequest(
		t, token, http.MethodPut,
		server.URL+"/users/narrator",
		bytes.NewBuffer(payload),
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}
}

func TestDeleteUser(t *testing.T) {
	server := Setup(t)

	token := loginUser(t, server, "admin", "admin")

	req := makeRequest(
		t, token, http.MethodDelete,
		server.URL+"/users/narrator",
		nil,
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("want 204, got %d", resp.StatusCode)
	}
}
