package articles_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/articles"
	"ssb/internal/domain/models"
	"ssb/internal/dto"
	"testing"
)

type FakeArticleRepository struct {
	Store map[string]models.Article
}

func (f *FakeArticleRepository) GetByID (id string) (models.Article, error){
	return models.Article{}, errors.New("Not Implimented")
}


func (f *FakeArticleRepository) ListAll()([]models.Article, error){
	return []models.Article{}, errors.New("Not Implimented")
}

func (f *FakeArticleRepository) Create(a dto.ArticleCreateDTO) (uint32, error){
	return 0, errors.New("Not Implimented")
}

func (f *FakeArticleRepository) Update(id string, update dto.ArticleUpdateDTO) error {
	return errors.New("Not Implimented")
}

func (f *FakeArticleRepository) Delete(id string)error {
	return errors.New("Not Implimented")
}

func NewFakeArticleRepository() *FakeArticleRepository{
	f := FakeArticleRepository{
		Store: make(map[string]models.Article),
	}
	return &f
}

func TestGetArticles(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r := articles.NewRouter(NewFakeArticleRepository())
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	expectedBody := "\"\"\n"

	if w.Body.String() != expectedBody {
		t.Errorf("unexpected body: got %q, want %q", w.Body.String(), expectedBody)
	}
}
