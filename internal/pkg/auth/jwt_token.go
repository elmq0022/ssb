package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"ssb/internal/schemas"
	"ssb/internal/timeutil"
	"time"
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

func (c *JWTConfig) DecodeToken(jsonToken schemas.JsonToken) (*jwt.RegisteredClaims, bool) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		jsonToken.Token,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(c.Secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		log.Fatal(err)
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	return claims, ok && token.Valid

}

func (c *JWTConfig) IsValidToken(
	username string,
	jsonToken schemas.JsonToken,
) (bool, error) {
	claims, ok := c.DecodeToken(jsonToken)
	if !ok || claims == nil {
		return false, errors.New("bad token")
	}

	// Validate subject
	if claims.Subject != username {
		return false, errors.New("subject mismatch")
	}

	// Validate issuer
	if claims.Issuer != c.Iss {
		return false, errors.New("issuer mismatch")
	}

	// Validate audience
	if len(claims.Audience) == 0 || claims.Audience[0] != c.Aud {
		return false, errors.New("audience mismatch")
	}

	// Validate timestamps
	now := c.Clock.Now().UTC()
	if claims.IssuedAt != nil && claims.IssuedAt.Time.After(now) {
		return false, errors.New("issued at time is in the future")
	}
	if claims.NotBefore != nil && claims.NotBefore.Time.After(now) {
		return false, errors.New("not before time is in the future")
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(now) {
		return false, errors.New("token has expired")
	}

	return true, nil
}
