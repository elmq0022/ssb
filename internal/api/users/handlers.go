package users

import (
	"errors"
	"net/http"
	"ssb/internal/models"
	"ssb/internal/pkg/router"
	"ssb/internal/repo"
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
		return "", http.StatusNotImplemented, errors.New("NotImplemented")
	})

	r.Put("/{userName}", func(req *http.Request) (any, int, error) {
		return "", http.StatusNotImplemented, errors.New("NotImplemented")
	})

	r.Delete("/{userName}", func(req *http.Request) (any, int, error) {
		return "", http.StatusNotImplemented, errors.New("NotImplemented")
	})

	return r
}
