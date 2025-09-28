package articles_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/articles"
	"ssb/internal/models"
	"ssb/internal/schemas"
	"ssb/internal/testutil"
	"testing"

	// "github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func setup(
	t *testing.T,
	httpMethod string,
	url string,
	body io.Reader,
	as []models.Article,
	us []models.User,
	username string,
) (
	*httptest.ResponseRecorder,
	*testutil.FakeArticleRepository,
) {
	t.Helper()
	req := httptest.NewRequest(httpMethod, url, body)
	w := httptest.NewRecorder()
	ar := testutil.NewFakeArticleRepository(as, us)
	ur := testutil.NewFakeUserRepository(us)

	//TODO: pass the correct user
	auth := func(request *http.Request) (string, error) { return username, nil }
	r := articles.NewRouter(ar, ur, auth)
	r.ServeHTTP(w, req)
	return w, ar
}

func TestGetArticles(t *testing.T) {
	as := []models.Article{
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
	us := []models.User{
		{
			UserName:  "author0",
			FirstName: "aa",
			LastName:  "bb",
			Email:     "aa@example.com",
			CreatedAt: testutil.Fc0.FixedTime.Unix(),
			UpdatedAt: testutil.Fc0.FixedTime.Unix(),
		},
		{
			UserName:  "author1",
			FirstName: "cc",
			LastName:  "dd",
			Email:     "cc@example.com",
			CreatedAt: testutil.Fc0.FixedTime.Unix(),
			UpdatedAt: testutil.Fc0.FixedTime.Unix(),
		},
	}
	w, _ := setup(t, http.MethodGet, "/", nil, as, us, "author0")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var got []schemas.ArticleWithAuthorSchema
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal respone: %v", err)
	}

	// TODO: fix comparision
	//if !cmp.Equal(want, got) {
	//	t.Fatalf("%v", cmp.Diff(want, got))
	//}
}

func TestGetArticleByID(t *testing.T) {
	as := testutil.NewArticle(
		testutil.Fc0,
		testutil.WithID("0"),
		testutil.WithTitle("title0"),
		testutil.WithAuthor("author0"),
		testutil.WithBody("body0"),
	)
	us := models.User{
		UserName:  "author0",
		FirstName: "aa",
		LastName:  "bb",
		Email:     "aa@example.com",
		CreatedAt: testutil.Fc0.FixedTime.Unix(),
		UpdatedAt: testutil.Fc0.FixedTime.Unix(),
	}
	w, _ := setup(
		t,
		http.MethodGet,
		"/0",
		nil,
		[]models.Article{as},
		[]models.User{us},
		"author0",
	)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var got schemas.ArticleWithAuthorSchema
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal respone: %v", err)
	}

	//if !cmp.Equal(want, got) {
	//	t.Fatalf("%v", cmp.Diff(want, got))
	//}
}

func TestDeleteArticle(t *testing.T) {
	user := models.User{
		UserName:  "user1",
		FirstName: "robert",
		LastName:  "bobby",
	}
	article := testutil.NewArticle(
		testutil.Fc0,
		testutil.WithID("0"),
		testutil.WithAuthor(user.UserName),
	)
	w, ar := setup(
		t,
		http.MethodDelete,
		"/0",
		nil,
		[]models.Article{article},
		[]models.User{user},
		user.UserName,
	)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 OK, got %d", w.Code)
	}

	_, exists := ar.ArticleStore["0"]
	if exists {
		t.Fatalf("article with id '0' was not deleted from the store")
	}
}

func TestCreateArticle(t *testing.T) {
	user := models.User{
		UserName:  "user2",
		FirstName: "edward",
		LastName:  "norton",
	}
	newArticle := schemas.ArticleCreateSchema{
		Title: "title",
		Body:  "body",
	}
	data, err := json.Marshal(newArticle)
	if err != nil {
		t.Fatalf("could not marshal dto: %q", newArticle)
	}

	w, ar := setup(
		t,
		http.MethodPost,
		"/",
		bytes.NewBuffer(data),
		[]models.Article{},
		[]models.User{user},
		user.UserName,
	)

	if w.Code != http.StatusCreated {
		t.Fatalf("failed to post the article: %v", w.Code)
	}

	var Id string
	if err := json.Unmarshal(w.Body.Bytes(), &Id); err != nil {
		t.Fatalf("could not unmarshal newly created article: %q", err)
	}

	if ar.ArticleStore[Id].Author != user.UserName {
		t.Errorf("wanted author but got: %s", ar.ArticleStore[Id].Author)
	}

	if ar.ArticleStore[Id].Title != "title" {
		t.Errorf("wanted title but got: %s", ar.ArticleStore[Id].Title)
	}

	if ar.ArticleStore[Id].Body != "body" {
		t.Errorf("wanted body but got: %s", ar.ArticleStore[Id].Body)
	}
}

func TestUpdateArticle(t *testing.T) {
	tests := []struct {
		name  string
		title *string
		body  *string
		user  models.User
	}{
		{
			name:  "no updates",
			title: nil,
			body:  nil,
			user: models.User{
				UserName: "u1",
			},
		},
		{
			name:  "update title",
			title: &[]string{"new title"}[0],
			body:  nil,
			user: models.User{
				UserName: "u2",
			},
		},
		{
			name:  "update body",
			title: nil,
			body:  &[]string{"new body"}[0],
			user: models.User{
				UserName: "u3",
			},
		},
		{
			name:  "update all",
			title: &[]string{"new title"}[0],
			body:  &[]string{"new body"}[0],
			user: models.User{
				UserName: "u4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := uuid.New().String()
			want := schemas.ArticleUpdateSchema{
				Title: tt.title,
				Body:  tt.body,
			}
			endpoint := fmt.Sprintf("/%s", id)
			data, err := json.Marshal(want)
			if err != nil {
				t.Fatalf("could not marshal json for %v", want)
			}

			article := testutil.NewArticle(
				testutil.Fc0,
				testutil.WithID(id),
				testutil.WithAuthor(tt.user.UserName),
			)
			w, ar := setup(
				t,
				http.MethodPut,
				endpoint,
				bytes.NewBuffer(data),
				[]models.Article{article},
				[]models.User{tt.user},
				tt.user.UserName,
			)

			if w.Code != http.StatusOK {
				t.Fatalf(
					"expected status code: %d, but got status code: %d",
					http.StatusOK,
					w.Code,
				)
			}

			got := ar.ArticleStore[id]

			if want.Title != nil && *want.Title != got.Title {
				t.Errorf("want title: %s, got title: %s", *want.Title, got.Title)
			}

			if want.Body != nil && *want.Body != got.Body {
				t.Errorf("want title: %s, got title: %s", *want.Body, got.Body)
			}
		})
	}
}
