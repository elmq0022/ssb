package app

import (
	"log"
	"os"
)

type Config struct {
	DBPath string
	JWTSecret string
	AdminPassword string
	Port string
}


func LoadConfig() Config {
	return Config{}
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

