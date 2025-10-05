package app

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Port          string
	DBPath        string
	JWTSecret     string
	AdminPassword string
	JWTIssuer     string
	JWTAudience   string
	JWT_TTL       time.Duration
}

func LoadConfig() Config {
	return Config{
		Port:          getEnv("BFS_PORT", ":8080"),
		DBPath:        getEnv("BFS_DB_PATH", "./data/ssb.db"),
		JWTSecret:     mustEnv("BFS_AUTH_SECRET"),
		AdminPassword: getEnv("BFS_ADMIN_PASSWD", "admin"),
		JWTIssuer:     getEnv("BFS_JWT_ISSUER", "bfs"),
		JWTAudience:   getEnv("BFS_JWT_AUDIEN", "bfs"),
		JWT_TTL:       getEnvDuration("BFS_JWT_TTL", time.Hour),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return value
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if dur, err := time.ParseDuration(value); err == nil {
			return dur
		}
		log.Fatalf("invalidd duration for %s=%s", key, value)
	}
	return defaultValue
}
