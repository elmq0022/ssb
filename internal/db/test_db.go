package db

import (
	"database/sql"
	_ "embed"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

//go:embed schema.sql
var schema string

func MustNewTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Panicf("could not open db: %v", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		log.Panicf("could not enable foreign keys: %v", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		log.Panicf("could not create schema: %v", err)
	}

	return db
}
