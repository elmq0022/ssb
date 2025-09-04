package articles

import (
	"encoding/json"
	"net/http"
	"ssb/internal/pkg/router"
	"ssb/internal/repo"
	"ssb/internal/schemas"
)

func NewRouter(ar repo.ArticleRepository) *router.Router {
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

	// must be authenticated
	r.Delete("/{id}", func(req *http.Request) (any, int, error) {
		id := req.PathValue("id")
		err := ar.Delete(id)
		if err != nil {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusOK, nil
	})

	// must be authenticated
	r.Post("/", func(req *http.Request) (any, int, error) {
		var data schemas.ArticleCreateSchema
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			return nil, http.StatusBadRequest, err
		}
		article, err := ar.Create(data)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return article, http.StatusCreated, nil
	})

	// must be authenticated
	r.Put("/{id}", func(req *http.Request) (any, int, error) {
		var update schemas.ArticleUpdateSchema
		Id := req.PathValue("id")
		if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
			return nil, http.StatusBadRequest, err
		}
		if err := ar.Update(Id, update); err != nil {
			return nil, http.StatusBadRequest, err
		}
		return nil, http.StatusOK, nil
	})

	return r
}
