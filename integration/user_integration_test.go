//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"ssb/internal/schemas"
	"testing"
	"ssb/integration/testutil"
)

func TestGetUser(t *testing.T) {
	server, _, _ := Setup(t)

	resp, err := http.Get(server.URL + "/users/admin")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d, but %d", http.StatusOK, resp.StatusCode)
	}
}

func TestCreateUser(t *testing.T) {
	server, _, _ := Setup(t)

	token := testutil.LoginUser(t, server, "admin", "admin")

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
	server, _, _ := Setup(t)

	token := testutil.LoginUser(t, server, "admin", "admin")

	updatedEmail := "narrator@paperstreetsoap.com"
	updateUserData := schemas.UpdateUserDTO{
		Email: &updatedEmail,
	}
	payload, err := json.Marshal(updateUserData)
	if err != nil {
		t.Fatalf("%v", err)
	}

	req := testutil.MakeAuthorizedRequest(
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
	server, _, _ := Setup(t)

	token := testutil.LoginUser(t, server, "admin", "admin")

	req := testutil.MakeAuthorizedRequest(
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
