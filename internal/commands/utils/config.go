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
		log.Fatalf("could not unmarshal config data: %v", data)
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
		log.Fatalf("could not unmarshal token data: %v", data)
	}
	return token
}

func ensureDir(path string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		log.Fatalf("could not create directory %q: %v", dir, err)
	}
}

func MustSetJWTToken(token schemas.JsonToken) {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		log.Fatalf("could not marshal jwt token due to err: %q", err)
	}
	ensureDir(JWTFilePath)
	os.WriteFile(JWTFilePath, data, 0o600)
}

func MustSetConfig(cfg *CLIConfig){
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatalf("could not marshal config err: %q", err)
	}
	ensureDir(ConfigFilePath)
	os.WriteFile(ConfigFilePath, data, 0o600)
}
