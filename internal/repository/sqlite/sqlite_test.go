package repo_test

import (
	"database/sql"
	"fmt"
	"log"
	"ssb/internal/domain/models"
	"ssb/internal/repository/sqlite"
	"ssb/internal/testutil"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestImports(t *testing.T) {
	r := repo.SqliteArticleRepo{}
	fmt.Printf("%v", r)
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
        id INTEGER UNIQUE,
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
	db := NewTestDB()
	want := testutil.DefaultArticle()

	db.Exec(INSERT_ARTICLE,
		want.ID,
		want.Title,
		want.Author,
		want.Body,
		want.PublishedAt.UTC().Format(time.RFC3339Nano),
		want.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)

	r := repo.NewSqliteArticleRepo(db)
	got, err := r.GetByID(want.ID)

	if err != nil {
		t.Fatal(err)
	}

	testutil.AssertArticleEqual(t, got, want)
}

func TestGetAllArticles(t *testing.T) {
	db := NewTestDB()

	a1 := testutil.DefaultArticle()
	a1.ID = 1

	a2 := testutil.DefaultArticle()
	a2.ID = 2
	a2.Title = "Article 2"
	a2.Body = "Body 2"
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

	r := repo.NewSqliteArticleRepo(db)
	got, err := r.ListAll()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !cmp.Equal(want, got) {
		t.Errorf("mismatch (-want +got):\n%s", cmp.Diff(want, got))
	}
}
