package main

import (
	"log"
	"net/http"
	"ssb/internal/api/articles"
	"ssb/internal/api/healthz"
	"ssb/internal/api/users"
	appDB "ssb/internal/db"
	"ssb/internal/pkg/auth"
	"ssb/internal/pkg/router"
	"ssb/internal/repo/sqlite"
	"ssb/internal/timeutil"
	"time"

	"github.com/jmoiron/sqlx"
)

func getJWTConfig() *auth.JWTConfig {
	config := auth.NewJWTConfig(
		auth.WithIssuer("ssb"),
		auth.WithAudience("ssb"),
		auth.WithTTL(time.Duration(1*time.Hour)),
		auth.WithClock(timeutil.RealClock{}),
		auth.WithSecretFromEnv("AUTH_SECRET"),
	)
	return config
}

func getOrCreateDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		log.Panicf("could not open db: %v", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		log.Panicf("could not enable foreign keys: %v", err)
	}

	if _, err := db.Exec(appDB.Schema); err != nil {
		db.Close()
		log.Panicf("could not create schema: %v", err)
	}

	return db
}

func main() {
	db := getOrCreateDB()
	clock := timeutil.RealClock{}
	ar := repo.NewSqliteArticleRepo(db, clock)
	ur := repo.NewUserSqliteRepo(db, clock)

	config := getJWTConfig()
	jwtAuth := router.NewJWTAuthFunction(config)

	mux := router.NewRouter()
	mux.Mount("/healthz", healthz.NewRouter())
	mux.Mount("/users", users.NewRouter(ur, jwtAuth))
	mux.Mount("/articles", articles.NewRouter(ar, ur, jwtAuth))
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
