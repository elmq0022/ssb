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

func NewRouter(ur repo.UserRepository, authFunc router.AuthFunc) *router.Router {
	r := router.NewRouter()

	r.Get("/{userName}", func(req *http.Request) (any, int, error) {
		userName := req.PathValue("userName")
		user, err := ur.GetByUserName(userName)

		if err != nil {
			return models.User{}, http.StatusNotFound, err
		}
		return user, http.StatusOK, nil
	})

	post := func(req *http.Request) (any, int, error) {
		user, ok := router.UserFromContext(req.Context())
		if !ok {
			return nil, http.StatusUnauthorized, errors.New("no username in context")
		}
		if user.UserName != "admin" {
			return nil, http.StatusUnauthorized, errors.New("permission denied")
		}

		var data schemas.CreateUserDTO
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return "", http.StatusBadRequest, err
		}
		userName, err := ur.Create(data)
		if err != nil {
			return "", http.StatusBadRequest, err
		}
		return userName, http.StatusCreated, nil
	}
	r.Post("/", router.WithAuth(post, authFunc, ur))

	put := func(req *http.Request) (any, int, error) {
		user, ok := router.UserFromContext(req.Context())
		if !ok {
			return nil, http.StatusUnauthorized, errors.New("no username in context")
		}
		if user.UserName != "admin" {
			return nil, http.StatusUnauthorized, errors.New("permission denied")
		}

		var data schemas.UpdateUserDTO
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return "", http.StatusBadRequest, err
		}
		userName := req.PathValue("userName")
		if err := ur.Update(userName, data); err != nil {
			return "", http.StatusBadRequest, err
		}
		return "", http.StatusOK, nil
	}
	r.Put("/{userName}", router.WithAuth(put, authFunc, ur))

	rm := func(req *http.Request) (any, int, error) {
		user, ok := router.UserFromContext(req.Context())
		if !ok {
			return nil, http.StatusUnauthorized, errors.New("no username in context")
		}
		if user.UserName != "admin" {
			return nil, http.StatusUnauthorized, errors.New("permission denied")
		}
		userName := req.PathValue("userName")
		if err := ur.Delete(userName); err != nil {
			return "", http.StatusBadRequest, errors.New("bad request")
		}
		return "", http.StatusNoContent, nil
	}
	r.Delete("/{userName}", router.WithAuth(rm, authFunc, ur))

	return r
}
