//go:build integration
// +build integration

package integration

import (
	// "net/http"
	// "ssb/integration/testutil"
	"ssb/internal/repo/sqlite"
	"ssb/internal/schemas"
	"testing"
)

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

	var a []string

	id1, _ := ar.Create(u1.UserName, a1)
	a = append(a, id1)

	id2, _ := ar.Create(u2.UserName, a2)
	a = append(a, id2)

	return a
}

func TestGetArticles(t *testing.T) {
	//server, ur, ar := Setup(t)
	//articleIds := createArticlesAndUsers(t, ur, ar)

	// list articles
	//testutil.MakeRequest(t, http.MethodGet, )

	// get article 1

	// get article 2


	t.Fatal("failed")
}
