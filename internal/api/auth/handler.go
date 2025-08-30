package auth

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"ssb/internal/pkg/auth"
	"ssb/internal/pkg/router"
	"ssb/internal/repo"
	"ssb/internal/schemas"
	"time"
)

func NewRouter(ur repo.UserRepository) *router.Router {
	r := router.NewRouter()

	r.Post("/login", func(req *http.Request) (any, int, error) {
		var data schemas.LoginRequest
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return "", http.StatusBadRequest, err
		}

		username := data.Username
		password := data.Password

		u, err := ur.GetByUserName(username)
		if err != nil {
			return "", http.StatusBadRequest, err
		}

		match, err := auth.CheckPassword(password, u.HashedPassword)
		if match && u.IsActive {
			// TODO: but in its own function
			//ref: https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-New-Hmac
			now := time.Now().UTC()
			exp := now.Add(1 * time.Hour)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": username,
				"iss": "ssb-auth-service",
				"aud": "ssb-backend",
				"iat": now.Unix(),
				"nbf": now.Unix(),
				"exp": exp.Unix(),
			})

			tokenString, err := token.SignedString([]byte("sample-secret"))
			if err != nil {
				return "", http.StatusInternalServerError, err
			}

			jwtToken := schemas.JsonToken{
				Token: tokenString,
			}

			return jwtToken, http.StatusOK, nil
		}

		return "", http.StatusUnauthorized, errors.New("Bad User")
	})

	return r
}
