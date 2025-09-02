package auth

import (
	"errors"
	"os"
	"ssb/internal/schemas"
	"ssb/internal/timeutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Iss    string
	Aud    string
	TTL    time.Duration
	Clock  timeutil.Clock
	Secret string
}

type JWTOption func(*JWTConfig)

func WithIssuer(iss string) JWTOption {
	return func(config *JWTConfig) {
		config.Iss = iss
	}
}

func WithAudience(aud string) JWTOption {
	return func(config *JWTConfig) {
		config.Aud = aud
	}
}

func WithTTL(ttl time.Duration) JWTOption {
	return func(config *JWTConfig) {
		config.TTL = ttl
	}
}
func WithClock(clock timeutil.Clock) JWTOption {
	return func(c *JWTConfig) {
		c.Clock = clock
	}
}

func WithSecret(secret string) JWTOption {
	return func(c *JWTConfig) {
		c.Secret = secret
	}
}

func WithSecretFromEnv(envName string) JWTOption {
	return func(c *JWTConfig) {
		c.Secret = os.Getenv(envName)
	}
}

func NewJWTConfig(opts ...JWTOption) *JWTConfig {
	c := &JWTConfig{
		TTL:   1 * time.Hour,
		Clock: timeutil.RealClock{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *JWTConfig) GenerateJWT(
	username string,
) (schemas.JsonToken, error) {
	now := c.Clock.Now().UTC()
	exp := now.Add(c.TTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"iss": c.Iss,
		"aud": c.Aud,
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": exp.Unix(),
	})
	if c.Secret == "" {
		return schemas.JsonToken{}, errors.New("no secret provided")
	}
	tokenString, err := token.SignedString([]byte(c.Secret))
	if err != nil {
		return schemas.JsonToken{}, err
	}
	jwtToken := schemas.JsonToken{
		Token: tokenString,
	}
	return jwtToken, nil
}
