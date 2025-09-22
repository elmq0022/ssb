package router_test

import (
	"net/http"
	"ssb/internal/pkg/auth"
	"ssb/internal/pkg/router"
	"testing"
)

func TestJwtAuthFunc(t *testing.T) {
	jwtConf := auth.NewJWTConfig(
		auth.WithAudience("ssb"),
		auth.WithIssuer("ssb"),
		auth.WithSecret("123"),
	)

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	jsonToken, err := jwtConf.GenerateJWT("username")
	if err != nil {
		t.Fatalf("%v", err)
	}
	req.Header.Set("Authorization", "Bearer "+jsonToken.Token)
	req.Header.Set("Content-Type", "application/json")

	authFunc := router.NewJWTAuthFunction(jwtConf)
	subject, err := authFunc(req)

	if subject != "username" {
		t.Fatalf("want username, got %s", subject)
	}
}
