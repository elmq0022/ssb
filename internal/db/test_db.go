package db

import (
	"database/sql"
	"embed"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

func NewTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
