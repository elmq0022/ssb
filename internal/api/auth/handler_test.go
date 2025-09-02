package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/auth"
	"ssb/internal/models"
	authutil "ssb/internal/pkg/auth"
	"ssb/internal/schemas"
	"ssb/internal/testutil"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	username := "bud.bill"
	password := "password"

	body := schemas.LoginRequest{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("could not marshal password and body for %v", data)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	userData := schemas.CreateUserDTO{
		UserName:  username,
		FirstName: "bud",
		LastName:  "bill",
		Email:     "bud.bill@kill.com",
		Password:  password,
	}

	ur := testutil.NewFakeUserRepository([]models.User{})
	ur.Create(userData)

	c := authutil.NewJWTConfig(
		authutil.WithAudience("ssb"),
		authutil.WithSecret("password"),
	)
	r := auth.NewRouter(ur, c)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want %d, got %d", http.StatusOK, w.Code)
	}

	var j schemas.JsonToken
	if err := json.NewDecoder(w.Body).Decode(&j); err != nil {
		t.Fatalf("bad marshal: %v", err)
	}

	// TODO: check the JWT returned.
	if j.Token == "" {
		t.Errorf("did not want an empty string '%s'", w.Body.String())
	}
}
