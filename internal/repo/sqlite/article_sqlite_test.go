package repo_test

import (
	tdb "ssb/internal/db"
	"ssb/internal/models"
	"ssb/internal/repo/sqlite"
	"ssb/internal/schemas"
	"ssb/internal/testutil"
	"ssb/internal/timeutil"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

// TODO move to test utils
func asserEqual(t *testing.T, want, got any) {
	t.Helper()
	if !cmp.Equal(want, got) {
		t.Errorf("mismatch (-want +got):\n%s", cmp.Diff(want, got))
	}
}

const INSERT_ARTICLE = `
INSERT INTO articles (
	id,
	title,
	author,
	body,
	published_at,
	updated_at
) VALUES (
	:id,
	:title,
	:author,
	:body,
	:published_at,
	:updated_at
)`

const INSERT_USER = `
INSERT INTO users (
  user_name,
  first_name,
  last_name,
  email,
  hashed_password,
  is_active,
  created_at,
  updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

func insertArticle(t *testing.T, db *sqlx.DB, article models.Article) {
	t.Helper()
	_, err := db.NamedExec(INSERT_ARTICLE, article)
	if err != nil {
		t.Fatalf("could not insert article into db: %v", err)
	}
}

func NewTestRepo(clock timeutil.Clock) (repo.SqliteArticleRepo, *sqlx.DB) {
	db := tdb.MustNewTestDB()
	r := repo.NewSqliteArticleRepo(db, clock)
	return r, db
}

func createUser(
	t *testing.T,
	db *sqlx.DB,
	userName,
	firstName,
	lastName,
	email string) {
	t.Helper()
	_, err := db.Exec(
		INSERT_USER,
		userName,
		firstName,
		lastName,
		email,
		"super_secret",
		true,
		testutil.Fc0.FixedTime.Unix(),
		testutil.Fc0.FixedTime.Unix(),
	)

	if err != nil {
		t.Fatalf("could not create user: %v", err)
	}
}

func TestGetArticleByID(t *testing.T) {
	r, db := NewTestRepo(testutil.Fc0)
	t.Cleanup(func() { db.Close() })

	want := testutil.NewArticle(testutil.Fc0)
	createUser(t, db, want.Author, "FirstName", "LastName", "email@example.com")
	insertArticle(t, db, want)

	got, err := r.GetByID(want.ID)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if got.Author.UserName != want.Author {
		t.Errorf("want: %v, got: %v", want.Author, got.Author.UserName)
	}

	if want.Title != got.Title {
		t.Errorf("want: %s, got: %s", want.Title, got.Title)
	}

	if want.Body != got.Body {
		t.Errorf("want: %s, got: %s", want.Body, got.Body)
	}
}

func TestGetAllArticles(t *testing.T) {
	r, db := NewTestRepo(testutil.Fc0)
	t.Cleanup(func() { db.Close() })

	userName1 := "martin"
	firstName1 := "m"
	lastName1 := "f"
	email1 := "f.m@radio.com"
	userName2 := "folwer"
	firstName2 := "n"
	lastName2 := "e"
	email2 := "n.e@one.com"

	createUser(t, db, userName1, firstName1, lastName1, email1)
	createUser(t, db, userName2, firstName2, lastName2, email2)

	a1 := testutil.NewArticle(
		testutil.Fc0,
		testutil.WithAuthor(userName1),
	)
	a2 := testutil.NewArticle(
		testutil.Fc0,
		testutil.WithID("2"),
		testutil.WithAuthor(userName2),
		testutil.WithTitle("Title 2"),
		testutil.WithBody("Body 2"),
	)
	insertArticle(t, db, a1)
	insertArticle(t, db, a2)

	want := []schemas.ArticleWithAuthorSchema{
		{
			Title: a1.Title,
			Body:  a1.Body,
			Author: schemas.UserBrief{
				UserName:  userName1,
				FirstName: firstName1,
				LastName:  lastName1,
			},
		},
		{
			Title: a2.Title,
			Body:  a2.Body,
			Author: schemas.UserBrief{
				UserName:  userName2,
				FirstName: firstName2,
				LastName:  lastName2,
			},
		},
	}

	got, err := r.ListAll()
	if err != nil {
		t.Fatalf("%v", err)
	}

	asserEqual(t, want, got)
}

func TestCreateArticle(t *testing.T) {
	r, db := NewTestRepo(testutil.Fc0)
	t.Cleanup(func() { db.Close() })

	title := "New Title"
	author := "New Author"
	body := "New Body"

	createUser(t, db, author, "first", "last", "email")
	create_dto := schemas.ArticleCreateSchema{
		Title:    title,
		UserName: author,
		Body:     body,
	}

	id, err := r.Create(create_dto)
	if err != nil {
		t.Fatalf("%q", err)
	}

	want := testutil.NewArticle(
		testutil.Fc0,
		testutil.WithID(id),
		testutil.WithTitle(title),
		testutil.WithAuthor(author),
		testutil.WithBody(body),
	)

	q := `SELECT
	  id,
	  title,
	  author,
	  body,
	  published_at,
	  updated_at
	FROM articles
	WHERE id=$1`

	var got = models.Article{}
	if err := db.Get(&got, q, id); err != nil {
		t.Fatalf("%v", err)
	}

	asserEqual(t, want, got)
}

func ptrFromString(s string) *string {
	return &s
}

// TODO: Need to populated DB with an article.
// Then do the update and then check the result.

func TestUpdateArticle(t *testing.T) {
	tests := []struct {
		name    string
		a       models.Article
		updates schemas.ArticleUpdateSchema
		want    models.Article
	}{
		{
			"no-op",
			testutil.NewArticle(testutil.Fc0),
			schemas.ArticleUpdateSchema{Title: nil, Body: nil},
			testutil.NewArticle(
				testutil.Fc0,
				testutil.WithUpdatedAt(testutil.Fc5)),
		},
		{
			"update-all",
			testutil.NewArticle(testutil.Fc0),
			schemas.ArticleUpdateSchema{
				Title:    ptrFromString("newTitle"),
				Body:     ptrFromString("newBody"),
			},
			testutil.NewArticle(
				testutil.Fc0,
				testutil.WithTitle("newTitle"),
				testutil.WithBody("newBody"),
				testutil.WithUpdatedAt(testutil.Fc5)),
		},
		{
			"update-title",
			testutil.NewArticle(testutil.Fc0),
			schemas.ArticleUpdateSchema{
				Title:    ptrFromString("newTitle"),
				Body:     nil,
			},
			testutil.NewArticle(
				testutil.Fc0,
				testutil.WithTitle("newTitle"),
				testutil.WithUpdatedAt(testutil.Fc5),
			),
		},
		{
			"update-body",
			testutil.NewArticle(testutil.Fc0),
			schemas.ArticleUpdateSchema{
				Title:    nil,
				Body:     ptrFromString("newBody"),
			},
			testutil.NewArticle(
				testutil.Fc0,
				testutil.WithBody("newBody"),
				testutil.WithUpdatedAt(testutil.Fc5),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, db := NewTestRepo(testutil.Fc5)
			
			createUser(t, db, tt.a.Author, "firstName", "lastName", "email@example.com")

			t.Cleanup(func() { db.Close() })

			sql := `INSERT INTO articles (
				id,
				title,
				author,
				body,
				published_at,
				updated_at
			)
			VALUES (?, ?, ?, ?, ?, ?)`

			_, err := db.Exec(sql,
				tt.a.ID,
				tt.a.Title,
				tt.a.Author,
				tt.a.Body,
				tt.a.PublishedAt.Format(time.RFC3339Nano),
				tt.a.UpdatedAt.Format(time.RFC3339Nano),
			)

			if err != nil {
				t.Fatalf("Could not create update target in db. Error: %q", err)
			}

			if err := r.Update(tt.a.ID, tt.updates); err != nil {
				t.Errorf("could not update article %v, error: %v", tt.a, err)
			}

			q := `SELECT
			  id,
			  title,
			  author,
			  body, 
			  published_at,
			  updated_at
			FROM articles
			WHERE id=$1`
			
			var got = models.Article{}
			if err := db.Get(&got, q, tt.a.ID); err != nil {
				t.Fatalf("%v", err)
			}
			asserEqual(t, tt.want, got)
		})
	}
}

func TestDeleteArticle(t *testing.T) {
	r, db := NewTestRepo(testutil.Fc0)
	t.Cleanup(func() { db.Close() })
	id := "delete-me"
	a := testutil.NewArticle(testutil.Fc0, testutil.WithID(id))

	createUser(t, db, a.Author, "first", "last", "email@example.com")
	insertArticle(t, db, a)


	if err := r.Delete(id); err != nil {
		t.Fatalf("%q", err)
	}

	want := 0
	check := `SELECT COUNT(*) FROM articles WHERE id = ?`

	var got int
	if err := db.QueryRow(check, a.ID).Scan(&got); err != nil {
		t.Fatalf("%q", err)
	}
	if want != got {
		t.Errorf("want: %d, got %d", want, got)
	}
}
