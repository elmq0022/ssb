package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"ssb/internal/schemas"
)

type JWT interface {
	GenerateJWT(username string) (schemas.JsonToken, error)
	DecodeTokenString(token string) (*jwt.RegisteredClaims, error)
	IsValidToken(username string) (*jwt.RegisteredClaims, error)
}
