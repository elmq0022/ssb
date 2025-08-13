package articles_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
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
	} else {
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

func (f *FakeArticleRepository) Create(a dto.ArticleCreateDTO) (string, error) {
	id := uuid.New().String()
	article := testutil.NewArticle(
		testutil.Fc0,
		testutil.WithID(id),
		testutil.WithAuthor(a.Author),
		testutil.WithTitle(a.Title),
		testutil.WithBody(a.Body),
	)
	f.Store[id] = article
	return id, nil
}

func (f *FakeArticleRepository) Update(id string, update dto.ArticleUpdateDTO) error {
	article, ok := f.Store[id]

	if !ok {
		return errors.New("article not found")
	}

	if update.Title != nil {
		article.Title = *update.Title
	}

	if update.Author != nil {
		article.Author = *update.Author
	}

	if update.Body != nil {
		article.Body = *update.Body
	}

	return nil
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

func TestCreateArticle(t *testing.T) {
	newArticle := dto.ArticleCreateDTO{
		Author: "author",
		Title:  "title",
		Body:   "body",
	}

	data, err := json.Marshal(newArticle)
	if err != nil {
		t.Fatalf("could not marshal dto: %q", newArticle)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(data))
	w := httptest.NewRecorder()

	ar := NewFakeArticleRepository([]models.Article{})
	r := articles.NewRouter(ar)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("failed to post the article: %v", w.Code)
	}

	var Id string
	if err := json.Unmarshal(w.Body.Bytes(), &Id); err != nil {
		t.Fatalf("could not unmarshal newly created article: %q", err)
	}

	if ar.Store[Id].Author != "author" {
		t.Errorf("wanted author but got: %s", ar.Store[Id].Author)
	}

	if ar.Store[Id].Title != "title" {
		t.Errorf("wanted title but got: %s", ar.Store[Id].Title)
	}

	if ar.Store[Id].Body != "body" {
		t.Errorf("wanted body but got: %s", ar.Store[Id].Body)
	}
}

func TestUpdateArticle(t *testing.T) {
	Id := uuid.New().String()
	endpoint := fmt.Sprintf("/%s", Id)
	want := dto.ArticleUpdateDTO{
		Title:  nil,
		Body:   nil,
		Author: nil,
	}

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("could not marshal json for %v", want)
	}

	req := httptest.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(data))

	article := testutil.NewArticle(testutil.Fc0, testutil.WithID(Id))
	ar := NewFakeArticleRepository([]models.Article{article})
	r := articles.NewRouter(ar)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf(
			"expected status code: %d, but got status code: %d",
			http.StatusOK,
			w.Code,
		)
	}

	got := ar.Store[Id]

	if want.Title != nil && *want.Title != got.Title {
		t.Errorf("want title: %s, got title: %s", *want.Title, got.Title)
	}

	if want.Author != nil && *want.Author != got.Author {
		t.Errorf("want title: %s, got title: %s", *want.Author, got.Author)
	}

	if want.Body != nil && *want.Body != got.Body {
		t.Errorf("want title: %s, got title: %s", *want.Body, got.Body)
	}
}
