package auth

import (
	"errors"
	"os"
	"slices"
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

	claims := jwt.RegisteredClaims{
		Subject:   username,
		Issuer:    c.Iss,
		Audience:  []string{c.Aud},
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(exp),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if c.Secret == "" {
		return schemas.JsonToken{}, errors.New("no secret provided")
	}

	tokenString, err := token.SignedString([]byte(c.Secret))
	if err != nil {
		return schemas.JsonToken{}, err
	}

	return schemas.JsonToken{Token: tokenString}, nil
}

func (c *JWTConfig) DecodeTokenString(tokenString string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(c.Secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid{
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (c *JWTConfig) IsValidToken(tokenString string) (*jwt.RegisteredClaims, error) {
	claims, err := c.DecodeTokenString(tokenString)
	if err != nil {
		return nil, err
	}

	now := c.Clock.Now().UTC()

	if claims.Issuer != c.Iss {
		return nil, errors.New("issuer mismatch")
	}

	if !slices.Contains(claims.Audience, c.Aud){
		return nil, errors.New("audience mismatch")
	}

	if claims.IssuedAt != nil && claims.IssuedAt.Time.After(now) {
		return nil, errors.New("issued at time is in the future")
	}
	if claims.NotBefore != nil && claims.NotBefore.Time.After(now) {
		return nil, errors.New("not before time is in the future")
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(now) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
