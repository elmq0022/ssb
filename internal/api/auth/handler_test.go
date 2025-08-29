package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/auth"
	"ssb/internal/schemas"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	body := schemas.LoginRequest{
		Username: "name",
		Password: "password",
	}

	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("could not marshal password and body for %v", data)
	}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := auth.NewRouter()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
	}
}
