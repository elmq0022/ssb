package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"ssb/internal/schemas"
)

var (
	ConfigFilePath string
	JWTFilePath    string
)

func init() {
	ConfigFilePath = MustGetConfigFile()
	JWTFilePath = MustGetJWTFile()
}

type CLIConfig struct {
	URL      string `json:"url"`
	Username string `json:"user"`
}

func MustGetConfigFile() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("could not get config dir: %q", err)
	}
	app := "bfs"
	file := "config.json"
	return filepath.Join(configDir, app, file)
}

func MustGetJWTFile() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("could not get cache dir: %q", err)
	}
	app := "bfs"
	file := "token.json"
	return filepath.Join(cacheDir, app, file)
}

func MustReadConfig() *CLIConfig {
	data, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatalf("couldn't read config: %q", err)
	}

	var c CLIConfig
	if err := json.Unmarshal(data, &c); err != nil {
		log.Fatalf("could not unmarshal config data: %s", data)
	}
	return &c
}

func MustReadJWTToken() schemas.JsonToken {
	data, err := os.ReadFile(JWTFilePath)
	if err != nil {
		log.Fatalf("couldn't read jwt: %q", err)
	}
	var token schemas.JsonToken
	if err := json.Unmarshal(data, &token); err != nil {
		log.Fatalf("could not unmarshal token data: %s", data)
	}
	return token
}

func MustSetJWTToken(token schemas.JsonToken) {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		log.Fatalf("could not marsal jwt token due to err: %q", err)
	}
	os.WriteFile(JWTFilePath, data, 0o600)
}
