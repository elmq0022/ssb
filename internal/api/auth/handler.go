package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"ssb/internal/pkg/auth"
	"ssb/internal/pkg/router"
	"ssb/internal/repo"
	"ssb/internal/schemas"
)

func NewRouter(ur repo.UserRepository, c auth.JWT) *router.Router {
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
			token, err := c.GenerateJWT(username)
			if err != nil {
				return schemas.JsonToken{}, http.StatusInternalServerError, err
			}
			return token, http.StatusOK, nil
		}
		return "", http.StatusUnauthorized, errors.New("Bad User")
	})

	return r
}
