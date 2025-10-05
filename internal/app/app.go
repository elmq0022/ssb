package app

import (
	"ssb/internal/db"
	"ssb/internal/schemas"
	"ssb/internal/repo/sqlite"
	"log"
	"os"
	"net/http"
	"ssb/internal/api/articles"
	authApi "ssb/internal/api/auth"
	"ssb/internal/api/healthz"
	"ssb/internal/api/users"
	"ssb/internal/pkg/router"
	"ssb/internal/pkg/auth"
	"ssb/internal/timeutil"
	"time"
	"github.com/jmoiron/sqlx"
	"context"
)

type App struct {
	DB *sqlx.DB
	Router http.Handler
	Shutdown chan struct{}
}

func NewApp() *App {
	db := db.GetOrCreateDB()
	clock := timeutil.RealClock{}
	ar := repo.NewSqliteArticleRepo(db, clock)
	ur := repo.NewUserSqliteRepo(db, clock)
	createAdmin(ur)

	config := getJWTConfig()
	jwtAuth := router.NewJWTAuthFunction(config)

	mux := router.NewRouter()
	mux.Mount("/healthz", healthz.NewRouter())
	mux.Mount("/users", users.NewRouter(ur, jwtAuth))
	mux.Mount("/articles", articles.NewRouter(ar, ur, jwtAuth))
	mux.Mount("/auth", authApi.NewRouter(ur, config))

    return &App{
        DB:       db,
        Router:   mux,
        Shutdown: make(chan struct{}),
    }
}

func (app *App) Run(){
	srv := &http.Server{
        Addr:    ":8080",
        Handler: app.Router,
    }

    go func() {
        log.Println("Server running on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server failed: %v", err)
        }
    }()

    <-app.Shutdown
    log.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
    app.DB.Close()
}


func createAdmin(ur *repo.UserSqliteRepo) {
	passwd := os.Getenv("BFS_ADMIN_PASSWD")
	if passwd == "" {
		passwd = "admin"
	}
	username := "admin"
	data := schemas.CreateUserDTO{
		UserName:  username,
		FirstName: "",
		LastName:  "",
		Password:  passwd,
	}
	_, err := ur.Create(data)
	if err != nil {
		log.Panic("could not create admin account")
	}
}

func getJWTConfig() *auth.JWTConfig {
	config := auth.NewJWTConfig(
		auth.WithIssuer("ssb"),
		auth.WithAudience("ssb"),
		auth.WithTTL(time.Duration(1*time.Hour)),
		auth.WithClock(timeutil.RealClock{}),
		auth.WithSecretFromEnv("BFS_AUTH_SECRET"),
	)
	return config
}
