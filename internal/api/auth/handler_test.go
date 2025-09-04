package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/auth"
	"ssb/internal/models"
	"ssb/internal/schemas"
	"ssb/internal/testutil"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

type FakeJWT struct {
	token string
}

func (f *FakeJWT) GenerateJWT(username string) (schemas.JsonToken, error) {
	return schemas.JsonToken{Token: f.token}, nil
}

func (f *FakeJWT) DecodeToken(JsonToken schemas.JsonToken) (*jwt.RegisteredClaims, bool) {
	return &jwt.RegisteredClaims{}, false
}

func (f *FakeJWT) IsValidToken(username string, JsonToken schemas.JsonToken) (bool, error) {
	return false, errors.New("not implemented")
}

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

	f := &FakeJWT{
		token: "valid-token",
	}
	r := auth.NewRouter(ur, f)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want %d, got %d", http.StatusOK, w.Code)
	}

	var j schemas.JsonToken
	if err := json.NewDecoder(w.Body).Decode(&j); err != nil {
		t.Fatalf("bad marshal: %v", err)
	}

	if j.Token != "valid-token" {
		t.Errorf("want: valid-token got: %s", j.Token)
	}
}
