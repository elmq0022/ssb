//go:build integration
// +build integration

package integration

import (
	"net/http/httptest"
	"ssb/internal/api/articles"
	authApi "ssb/internal/api/auth"
	"ssb/internal/api/healthz"
	"ssb/internal/api/users"
	appDB "ssb/internal/db"
	"ssb/internal/pkg/auth"
	"ssb/internal/pkg/router"
	"ssb/internal/repo/sqlite"
	"ssb/internal/schemas"
	"ssb/internal/timeutil"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func Setup(t *testing.T) (*httptest.Server, *repo.UserSqliteRepo, *repo.SqliteArticleRepo) {
	db := createDB(t)
	clock := timeutil.RealClock{}
	ar := repo.NewSqliteArticleRepo(db, clock)
	ur := repo.NewUserSqliteRepo(db, clock)
	createAdmin(t, ur)
	createUser(t, ur)

	config := getJWTConfig()
	jwtAuth := router.NewJWTAuthFunction(config)

	mux := router.NewRouter()
	mux.Mount("/healthz", healthz.NewRouter())
	mux.Mount("/users", users.NewRouter(ur, jwtAuth))
	mux.Mount("/articles", articles.NewRouter(ar, ur, jwtAuth))
	mux.Mount("/auth", authApi.NewRouter(ur, config))

	server := httptest.NewServer(mux)

	t.Cleanup(func() {
		db.Close()
		server.Close()
	})

	return server, ur, ar
}

func getJWTConfig() *auth.JWTConfig {
	config := auth.NewJWTConfig(
		auth.WithIssuer("ssb"),
		auth.WithAudience("ssb"),
		auth.WithTTL(time.Duration(1*time.Hour)),
		auth.WithClock(timeutil.RealClock{}),
		auth.WithSecret("testsecret"),
	)

	return config
}

func createAdmin(t *testing.T, ur *repo.UserSqliteRepo) {
	username := "admin"
	password := "admin"

	data := schemas.CreateUserDTO{
		UserName:  username,
		FirstName: "",
		LastName:  "",
		Password:  password,
	}

	_, err := ur.Create(data)
	if err != nil {
		t.Fatal("could not create admin account")
	}
}

func createUser(t *testing.T, ur *repo.UserSqliteRepo) {
	username := "narrator"
	password := "cornflower blue"

	data := schemas.CreateUserDTO{
		UserName:  username,
		FirstName: "edward",
		LastName:  "norton",
		Password:  password,
		Email:     "none",
	}

	_, err := ur.Create(data)
	if err != nil {
		t.Fatal("could not create narrator account")
	}
}

func createDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("could not open db: %v", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		t.Fatalf("could not enable foreign keys: %v", err)
	}

	if _, err := db.Exec(appDB.Schema); err != nil {
		db.Close()
		t.Fatalf("could not create schema: %v", err)
	}

	return db
}

