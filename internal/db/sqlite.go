package db

import (
	"github.com/jmoiron/sqlx"
	"log"
)

func GetOrCreateDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		log.Panicf("could not open db: %v", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		log.Panicf("could not enable foreign keys: %v", err)
	}

	if _, err := db.Exec(Schema); err != nil {
		db.Close()
		log.Panicf("could not create schema: %v", err)
	}

	return db
}
