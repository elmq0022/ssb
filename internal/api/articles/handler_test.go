package articles_test

import (
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/articles"
	"ssb/internal/domain/models"
	"ssb/internal/dto"
	"ssb/internal/testutil"
	"testing"
)

type FakeArticleRepository struct {
	Store map[string]models.Article
}

func (f *FakeArticleRepository) GetByID(id string) (models.Article, error) {
	article, exists := f.Store[id]
	if !exists {
		return models.Article{}, errors.New("Article Not Found")
	} else{
		return article, nil
	}
}

func (f *FakeArticleRepository) ListAll() ([]models.Article, error) {
	var articles []models.Article

	for _, v := range f.Store {
		articles = append(articles, v)
	}
	return articles, nil
}

func (f *FakeArticleRepository) Create(a dto.ArticleCreateDTO) (uint32, error) {
	return 0, errors.New("Not Implemented")
}

func (f *FakeArticleRepository) Update(id string, update dto.ArticleUpdateDTO) error {
	return errors.New("Not Implemented")
}

func (f *FakeArticleRepository) Delete(id string) error {
	_, exists := f.Store[id]
	if exists {
		delete(f.Store, id)
		return nil
	} else {
		return errors.New("Does not exist")
	}
}

func NewFakeArticleRepository(articles []models.Article) *FakeArticleRepository {
	s := make(map[string]models.Article)

	for _, article := range articles {
		s[article.ID] = article
	}

	f := FakeArticleRepository{
		Store: s,
	}

	return &f
}

func TestGetArticles(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	want := []models.Article{
		testutil.NewArticle(
			testutil.Fc0,
			testutil.WithID("0"),
			testutil.WithTitle("title0"),
			testutil.WithAuthor("author0"),
			testutil.WithBody("body0"),
		),
		testutil.NewArticle(
			testutil.Fc0,
			testutil.WithID("1"),
			testutil.WithTitle("title1"),
			testutil.WithAuthor("author1"),
			testutil.WithBody("body1"),
		),
	}

	r := articles.NewRouter(NewFakeArticleRepository(want))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var got []models.Article
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal respone: %v", err)
	}

	if !cmp.Equal(want, got) {
		t.Fatalf("%v", cmp.Diff(want, got))
	}
}

func TestGetArticleByID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/0", nil)
	w := httptest.NewRecorder()

	want := testutil.NewArticle(
		testutil.Fc0,
		testutil.WithID("0"),
		testutil.WithTitle("title0"),
		testutil.WithAuthor("author0"),
		testutil.WithBody("body0"),
	)

	r := articles.NewRouter(NewFakeArticleRepository([]models.Article{want}))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var got models.Article
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal respone: %v", err)
	}

	if !cmp.Equal(want, got) {
		t.Fatalf("%v", cmp.Diff(want, got))
	}
}

func TestDeleteArticle(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/0", nil)
	w := httptest.NewRecorder()

	article := testutil.NewArticle(
		testutil.Fc0,
		testutil.WithID("0"),
	)
	ar := NewFakeArticleRepository([]models.Article{article})
	r := articles.NewRouter(ar)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	_, exists := ar.Store["0"]
	if exists {
		t.Fatalf("article with id '0' was not deleted from the store")
	}
}
