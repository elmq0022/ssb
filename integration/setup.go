package integration

import (
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

func Setup(t *testing.T) *router.Router {
	db := createDB(t)
	clock := timeutil.RealClock{}
	ar := repo.NewSqliteArticleRepo(db, clock)
	ur := repo.NewUserSqliteRepo(db, clock)
	createAdmin(t, ur)

	config := getJWTConfig()
	jwtAuth := router.NewJWTAuthFunction(config)

	mux := router.NewRouter()
	mux.Mount("/healthz", healthz.NewRouter())
	mux.Mount("/users", users.NewRouter(ur, jwtAuth))
	mux.Mount("/articles", articles.NewRouter(ar, ur, jwtAuth))
	mux.Mount("/auth", authApi.NewRouter(ur, config))

	t.Cleanup(func() { db.Close() })

	return mux
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
	passwd := "admin"
	username := "admin"

	data := schemas.CreateUserDTO{
		UserName:  username,
		FirstName: "",
		LastName:  "",
		Password:  passwd,
	}

	_, err := ur.Create(data)
	if err != nil {
		t.Fatal("could not create admin account")
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
