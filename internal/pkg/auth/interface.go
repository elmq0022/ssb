package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"ssb/internal/schemas"
)

type JWT interface {
	GenerateJWT(username string) (schemas.JsonToken, error)
	DecodeToken(jsonToken schemas.JsonToken) (*jwt.RegisteredClaims, bool)
	IsValidToken(username string, jsonToken schemas.JsonToken) (bool, error)
}
