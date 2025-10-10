package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func MustGetConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("could not get config dir: %q", err)
	}
	app := "bfs"
	file := "config.json"
	return filepath.Join(configDir, app, file)
}

func MustGetCacheDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("could not get cache dir: %q", err)
	}
	app := "bfs"
	file := "token.json"
	return filepath.Join(cacheDir, app, file)
}

type CLIConfig struct {
	URL      string `json:"url"`
	Username string `json:"user"`
}

func LoadConfig() *CLIConfig {
	path := MustGetConfigDir()
	configData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("couldn't read config: %q", err)
	}

	var c CLIConfig
	if err := json.Unmarshal(configData, &c); err != nil {
		log.Fatalf("could not unmarshal config data: %s", configData)
	}
	return &c
}


func LoadJWTToken() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("could not get config dir: %q", err)
	}
	app := "bfs"
	file := "token.json"
	return filepath.Join(cacheDir, app, file)
}
