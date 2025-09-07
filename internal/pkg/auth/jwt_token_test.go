package auth_test

import (
	"ssb/internal/pkg/auth"
	"ssb/internal/timeutil"
	"testing"
	"time"
)

func GetTestJWTConfig() *auth.JWTConfig {
	clock := timeutil.FakeClock{
		FixedTime: time.Now().UTC(),
	}
	c := auth.NewJWTConfig(
		auth.WithClock(clock),
		auth.WithSecret("password"),
	)
	return c
}

func TestSuccessfulTokenGeneration(t *testing.T) {
	c := GetTestJWTConfig()
	token, err := c.GenerateJWT("username1")

	if err != nil {
		t.Fatalf("token generation error: %v", err)
	}

	if token.Token == "" {
		t.Fatal("token not valid")
	}
}

func TestSuccessfulTokenValidation(t *testing.T) {
	c := GetTestJWTConfig()
	encoded_token, err := c.GenerateJWT("username2")
	if err != nil {
		t.Fatalf("token generation error: %v", err)
	}

	claim, err := c.DecodeTokenString(encoded_token.Token)
	if err != nil {
		t.Fatalf("did not decode the token")
	}

	if claim.Subject != "username2" {
		t.Fatalf("want 'username2', got '%s'", claim.Subject)
	}
}

func TestTokenIsValid(t *testing.T) {
	c := GetTestJWTConfig()
	encoded_token, err := c.GenerateJWT("username2")
	if err != nil {
		t.Fatalf("token generation error: %v", err)
	}

	claims, err := c.IsValidToken(encoded_token.Token)
	if err != nil {
		t.Fatalf("error while checking token: %v", err)
	}
	if claims.Subject != "username2" {
		t.Fatalf("want username2, got %s", claims.Subject)
	}
}
