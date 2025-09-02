package auth

import "ssb/internal/schemas"

type JWT interface {
	GenerateJWT(username string) (schemas.JsonToken, error)
}
