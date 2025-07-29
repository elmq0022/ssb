package repo_test

import (
	"database/sql"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"log"
	"ssb/internal/domain/models"
	"ssb/internal/dto"
	"ssb/internal/repository/sqlite"
	"ssb/internal/testutil"
	"testing"
	"time"
)

func TestImports(t *testing.T) {
	r := repo.SqliteArticleRepo{}
	fmt.Printf("%v", r)
}

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
) VALUES (?, ?, ?, ?, ?, ?)`

func NewTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	schema := `
	    CREATE TABLE articles (
		pk INTEGER PRIMARY KEY AUTOINCREMENT,
        id TEXT UNIQUE NOT NULL,
        title TEXT NOT NULL,
        author TEXT NOT NULL,
        body TEXT NOT NULL,
        published_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    )`

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
	return db
}

func NewTestRepo() (repo.SqliteArticleRepo, *sql.DB) {
	db := NewTestDB()
	r := repo.NewSqliteArticleRepo(db, testutil.Fc0)
	return r, db
}

func TestNewDBReturnsZeroRows(t *testing.T) {
	db := NewTestDB()
	defer db.Close()

	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM articles").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatalf("wanted: 0 rows but got: %d rows", count)
	}
}

func TestGetArticleByID(t *testing.T) {
	r, db := NewTestRepo()
	defer db.Close()

	want := testutil.NewArticle(testutil.Fc0)
	db.Exec(INSERT_ARTICLE,
		want.ID,
		want.Title,
		want.Author,
		want.Body,
		want.PublishedAt.UTC().Format(time.RFC3339Nano),
		want.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)

	got, err := r.GetByID(want.ID)

	if err != nil {
		t.Fatal(err)
	}

	testutil.AssertArticleEqual(t, got, want)
}

func TestGetAllArticles(t *testing.T) {
	r, db := NewTestRepo()
	defer db.Close()

	a1 := testutil.NewArticle(testutil.Fc0)

	a2 := testutil.NewArticle(
		testutil.Fc5,
		testutil.WithID("2"),
		testutil.WithAuthor("Author 2"),
		testutil.WithTitle("Title 2"),
		testutil.WithBody("Body 2"),
	)

	want := []models.Article{a1, a2}
	for _, a := range want {
		db.Exec(INSERT_ARTICLE,
			a.ID,
			a.Title,
			a.Author,
			a.Body,
			a.PublishedAt.UTC().Format(time.RFC3339Nano),
			a.UpdatedAt.UTC().Format(time.RFC3339Nano),
		)
	}

	got, err := r.ListAll()
	if err != nil {
		t.Fatalf("%v", err)
	}

	asserEqual(t, want, got)
}

func TestCreateArticle(t *testing.T) {
	r, db := NewTestRepo()
	defer db.Close()

	title := "New Title"
	author := "New Author"
	body := "New Body"
	create_dto := dto.ArticleCreateDTO{
		Title:  title,
		Author: author,
		Body:   body,
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

	// TODO: I should probably use raw sql here.
	got, err := r.GetByID(id)
	asserEqual(t, want, got)
}

func ptrFromString(s string) *string {
	return &s
}

// TODO: Need to populated DB with an article.
// Then do the update and then check the result.
func TestUpdateArticle(t *testing.T) {
	r, db := NewTestRepo()
	defer db.Close()
	id := "10"
	update_dto := dto.ArticleUpdateDTO{
		Title:  ptrFromString("new title"),
		Author: ptrFromString("new author"),
		Body:   ptrFromString("new body"),
	}
	r.Update(id, update_dto)
	got, err := r.GetByID(id)
	if err != nil {
		t.Fatalf("%q", err)
	}

	fmt.Printf("%q", got)
}

