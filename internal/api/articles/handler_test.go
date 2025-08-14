package articles_test

import (
	"bytes"
	"encoding/json"
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

	r := articles.NewRouter(testutil.NewFakeArticleRepository(want))
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

	r := articles.NewRouter(testutil.NewFakeArticleRepository([]models.Article{want}))
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
	ar := testutil.NewFakeArticleRepository([]models.Article{article})
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

	ar := testutil.NewFakeArticleRepository([]models.Article{})
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
	tests := []struct {
		name   string
		author *string
		title  *string
		body   *string
	}{
		{
			name:   "no updates",
			author: nil,
			title:  nil,
			body:   nil,
		},
		{
			name:   "update author",
			author: &[]string{"new author"}[0],
			title:  nil,
			body:   nil,
		},
		{
			name:   "update title",
			author: nil,
			title:  &[]string{"new title"}[0],
			body:   nil,
		},
		{
			name:   "update body",
			author: nil,
			title:  nil,
			body:   &[]string{"new body"}[0],
		},
		{
			name:   "update all",
			author: &[]string{"new author"}[0],
			title:  &[]string{"new title"}[0],
			body:   &[]string{"new body"}[0],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Id := uuid.New().String()
			want := dto.ArticleUpdateDTO{
				Title:  tt.title,
				Author: tt.author,
				Body:   tt.body,
			}
			endpoint := fmt.Sprintf("/%s", Id)
			data, err := json.Marshal(want)
			if err != nil {
				t.Fatalf("could not marshal json for %v", want)
			}

			req := httptest.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(data))

			article := testutil.NewArticle(testutil.Fc0, testutil.WithID(Id))
			ar := testutil.NewFakeArticleRepository([]models.Article{article})
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
		})
	}
}
