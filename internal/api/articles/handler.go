package articles

import (
	"net/http"
	"ssb/internal/repository"
	"ssb/internal/router"
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

	return r
}
