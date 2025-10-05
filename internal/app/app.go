package app

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ssb/internal/api/articles"
	authApi "ssb/internal/api/auth"
	"ssb/internal/api/healthz"
	"ssb/internal/api/users"
	"ssb/internal/db"
	appDB "ssb/internal/db"
	"ssb/internal/pkg/auth"
	"ssb/internal/pkg/router"
	"ssb/internal/repo/sqlite"
	"ssb/internal/schemas"
	"ssb/internal/timeutil"
	"syscall"
	"time"
)

type App struct {
	Config Config
	DB     *sqlx.DB
	Server *http.Server
}

func NewApp(cfg Config) *App {
	database := db.OpenSQLite(cfg.DBPath, appDB.Schema)

	clock := timeutil.RealClock{}
	ar := repo.NewSqliteArticleRepo(database, clock)
	ur := repo.NewUserSqliteRepo(database, clock)

	createAdmin(ur, cfg.AdminPassword)

	jwtCfg := auth.NewJWTConfig(
		auth.WithIssuer(cfg.JWTIssuer),
		auth.WithAudience(cfg.JWTAudience),
		auth.WithTTL(cfg.JWT_TTL),
		auth.WithClock(clock),
		auth.WithSecret(cfg.JWTSecret),
	)
	jwtAuth := router.NewJWTAuthFunction(jwtCfg)

	mux := router.NewRouter()
	mux.Mount("/healthz", healthz.NewRouter())
	mux.Mount("/users", users.NewRouter(ur, jwtAuth))
	mux.Mount("/articles", articles.NewRouter(ar, ur, jwtAuth))
	mux.Mount("/auth", authApi.NewRouter(ur, jwtCfg))

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: mux,
	}

	return &App{
		Config: cfg,
		DB:     database,
		Server: srv,
	}
}

func (a *App) Run() error {
	log.Printf("Starting BFS server on %s\n", a.Config.Port)

	// Handle graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutdown signal received")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.Server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		a.DB.Close()
		close(idleConnsClosed)
	}()

	if err := a.Server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	log.Println("Server stopped gracefully")
	return nil
}

func createAdmin(ur *repo.UserSqliteRepo, password string) {
	username := "admin"
	data := schemas.CreateUserDTO{
		UserName:  username,
		FirstName: "",
		LastName:  "",
		Password:  password,
	}
	if _, err := ur.Create(data); err != nil {
		log.Panic("could not create admin account")
	}
}
