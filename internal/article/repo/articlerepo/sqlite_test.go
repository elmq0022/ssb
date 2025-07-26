package repo_test

import (
	"database/sql"
	"fmt"
	"log"
	"ssb/internal/article/repo/articlerepo"
	"ssb/internal/testutil"
	"testing"
	"time"
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
		want.Id,
		want.Title,
		want.Author,
		want.Body,
		want.PublishedAt.UTC().Format(time.RFC3339Nano),
		want.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)

	r := repo.NewSqliteArticleRepo(db)
	got, err := r.GetByID(want.Id)

	if err != nil {
		t.Fatal(err)
	}

	testutil.AssertArticleEqual(t, got, want)
}
