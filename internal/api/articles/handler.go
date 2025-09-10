package articles

import (
	"encoding/json"
	"errors"
	"net/http"
	"ssb/internal/pkg/router"
	"ssb/internal/repo"
	"ssb/internal/schemas"
)

func NewRouter(ar repo.ArticleRepository, ur repo.UserRepository, authFunc router.AuthFunc) *router.Router {
	r := router.NewRouter()

	r.Get("/", func(req *http.Request) (any, int, error) {
		articles, err := ar.ListAll()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return articles, http.StatusOK, nil
	})

	r.Get("/{id}", func(req *http.Request) (any, int, error) {
		id := req.PathValue("id")
		article, err := ar.GetByID(id)
		if err != nil {
			return nil, http.StatusNotFound, err
		}
		return article, http.StatusOK, nil
	})

	deleteHandler := func(req *http.Request) (any, int, error) {
		id := req.PathValue("id")
		err := ar.Delete(id)
		if err != nil {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusOK, nil
	}
	r.Delete("/{id}", router.WithAuth(deleteHandler, authFunc, ur))

	post := func(req *http.Request) (any, int, error) {
		var data schemas.ArticleCreateSchema
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return nil, http.StatusBadRequest, err
		}
		user, ok := router.UserFromContext(req.Context())
		if !ok {
			return nil, http.StatusUnauthorized, errors.New("no username in context")
		}
		article, err := ar.Create(user.UserName, data)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return article, http.StatusCreated, nil
	}
	r.Post("/", router.WithAuth(post, authFunc, ur))

	put := func(req *http.Request) (any, int, error) {
		var update schemas.ArticleUpdateSchema
		Id := req.PathValue("id")
		if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
			return nil, http.StatusBadRequest, err
		}
		if err := ar.Update(Id, update); err != nil {
			return nil, http.StatusBadRequest, err
		}
		return nil, http.StatusOK, nil
	}
	r.Put("/{id}", router.WithAuth(put, authFunc, ur))

	return r
}
