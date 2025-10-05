package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

func OpenSQLite(path string, schema string) *sqlx.DB {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Fatalf("couldl not create DB directory: %v", err)
	}

	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("could not open DB directory: %v", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		log.Panicf("could not enable foreign keys: %v", err)
	}

	ensureSchema(db, schema)
	return db
}

func ensureSchema(db *sqlx.DB, schema string) {
	var count int
	err := db.Get(
		&count,
		`SELECT count(*) FROM sqlite_master WHERE type='table' AND name='users';`,
	)
	if err == sql.ErrNoRows || count == 0 {
		log.Println("No schema found — applying schema...")
		if _, err := db.Exec(schema); err != nil {
			log.Fatalf("failed to apply schema: %v", err)
		}
	} else if err != nil {
		log.Fatalf("failed to check schema: %v", err)
	}
}
