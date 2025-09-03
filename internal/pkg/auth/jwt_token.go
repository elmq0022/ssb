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

func (c *JWTConfig) DecodeToken(jsonToken schemas.JsonToken) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(
		jsonToken.Token,
		func(token *jwt.Token) (any, error) {
			return []byte(c.Secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		log.Fatal(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}

func (c *JWTConfig) IsValidToken(
	username string,
	jsonToken schemas.JsonToken,
) (bool, error) {
	claims, ok := c.DecodeToken(jsonToken)
	if !ok {
		return false, errors.New("bad token")
	}
	// now := c.Clock.Now().UTC().Unix()
	if username != claims["sub"] {
		return false, nil
	}
	//TODO: test the rest of the claims
	return true, nil
}
