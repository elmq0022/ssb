//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"slices"
	"ssb/integration/testutil"
	"ssb/internal/repo/sqlite"
	"ssb/internal/schemas"
	"testing"
)

func testHttpClient(t *testing.T, req *http.Request)*http.Response{
	t.Helper()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	t.Cleanup(func(){resp.Body.Close()})
	return resp 
}

func createArticlesAndUsers(
	t *testing.T,
	ur *repo.UserSqliteRepo,
	ar *repo.SqliteArticleRepo,
) []string {
	t.Helper()

	u1 := schemas.CreateUserDTO{
		UserName: "user1",
		FirstName: "bob",
		LastName: "martin",
		Email: "bob@martin.com",
		Password: "secret1",
	}

	u2 := schemas.CreateUserDTO{
		UserName: "user2",
		FirstName: "david",
		LastName: "hanson",
		Email: "david@hanson.com",
		Password: "secret2",
	}

	ur.Create(u1)
	ur.Create(u2)

	a1 := schemas.ArticleCreateSchema{
		Title: "article1",
		Body: "this is article1",
	}

	a2 := schemas.ArticleCreateSchema{
		Title: "article2",
		Body: "this is article2",
	}

	var ids []string

	id1, _ := ar.Create(u1.UserName, a1)
	ids = append(ids, id1)

	id2, _ := ar.Create(u2.UserName, a2)
	ids = append(ids, id2)

	return ids
}

func TestGetArticles(t *testing.T) {
	server, ur, ar := Setup(t)
	articleIds := createArticlesAndUsers(t, ur, ar)

	// list articles
	req := testutil.MakeRequest(t, http.MethodGet, server.URL + "/articles", nil)
	resp := testHttpClient(t, req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}

	var articles []schemas.ArticleWithAuthorSchema
	json.NewDecoder(resp.Body).Decode(&articles)

	if  len(articles) != 2 {
		t.Fatalf("want 2 articles, got %d", len(articles))
	}

	for _, article := range articles{
		if  !slices.Contains(articleIds, article.ID) {
			t.Fatalf("returned articles does not contain title: %s", article.Title)
		}
	}
}

func TestGetArticleByID(t *testing.T) {
	server, ur, ar := Setup(t)
	articleIds := createArticlesAndUsers(t, ur, ar)
	articleId := articleIds[0]

	req := testutil.MakeRequest(t, http.MethodGet, server.URL + "/articles/" + articleId, nil)
	resp := testHttpClient(t, req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}
}

func TestUpdateArticle(t *testing.T){

}

func TestCreateArticle(t *testing.T){
	server, ur, ar := Setup(t)
	createArticlesAndUsers(t, ur, ar)
	token := testutil.LoginUser(t, server, "user2", "secret2")
	
	newArticleData := schemas.ArticleCreateSchema {
		Title: "New Title",
		Body: "New Body",
	}

	payload, err := json.Marshal(newArticleData)
	if err != nil {
		t.Fatalf("", )
	}

	req := testutil.MakeAuthorizedRequest(
		t,
		token,
		http.MethodPost,
		server.URL + "/articles",
		bytes.NewBuffer(payload),
	)

	resp := testHttpClient(t, req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf(
			"create article: want %d, got %d",
			http.StatusCreated,
			resp.StatusCode,
		)
	}
}

func TestDeleteArticle(t *testing.T){
	server, ur, ar := Setup(t)
	articleIDs := createArticlesAndUsers(t, ur, ar)
	articleID := articleIDs[1]

	token := testutil.LoginUser(t, server, "user2", "secret2")

	articleList, _ := ar.ListAll()
	articleCount := len(articleList)

	req := testutil.MakeAuthorizedRequest(
		t,
		token,
		http.MethodDelete,
		server.URL + "/articles/" + articleID,
		nil,
	)
	resp := testHttpClient(t, req)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("want %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	remainingArticles, _ := ar.ListAll()
	if articleCount - len(remainingArticles) != 1{
		t.Fatal("want 1 remaining article")
	}
}

