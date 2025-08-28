package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"ssb/internal/models"
	"ssb/internal/pkg/router"
	"ssb/internal/repo"
	"ssb/internal/schemas"
)

func NewRouter(ur repo.UserRepository) *router.Router {
	r := router.NewRouter()

	r.Get("/{userName}", func(req *http.Request) (any, int, error) {
		userName := req.PathValue("userName")
		user, err := ur.GetByUserName(userName)

		if err != nil {
			return models.User{}, http.StatusNotFound, err
		}
		return user, http.StatusOK, nil
	})

	r.Post("/", func(req *http.Request) (any, int, error) {
		var data schemas.CreateUserDTO
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return "", http.StatusBadRequest, err
		}
		userName, err := ur.Create(data)
		if err != nil {
			return "", http.StatusBadRequest, err
		}
		return userName, http.StatusCreated, nil
	})

	r.Put("/{userName}", func(req *http.Request) (any, int, error) {
		return "", http.StatusNotImplemented, errors.New("NotImplemented")
	})

	r.Delete("/{userName}", func(req *http.Request) (any, int, error) {
		userName := req.PathValue("userName")
		if err := ur.Delete(userName); err != nil {
			return "", http.StatusBadRequest, errors.New("bad request")
		}
		return "", http.StatusNoContent, nil
	})

	return r
}
