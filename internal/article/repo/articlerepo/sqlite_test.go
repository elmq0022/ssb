package repo_test

import (
	"database/sql"
	"fmt"
	"log"
	"ssb/internal/article/repo/articlerepo"
	"testing"
)

func TestImports(t *testing.T) {
	r := repo.SqliteArticleRepo{}
	fmt.Printf("%v", r)
}

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
