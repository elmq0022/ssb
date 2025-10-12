package utils_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"ssb/internal/commands/utils"
	"ssb/internal/schemas"
	"testing"
)

func TestMustGetConfigFile(t *testing.T) {
	home, _ := os.UserHomeDir()
	configDir := ".config"
	app := "bfs"
	file := "config.json"
	want := filepath.Join(home, configDir, app, file)
	got := utils.MustGetConfigFile()
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestMustGetJWTFile(t *testing.T) {
	home, _ := os.UserHomeDir()
	cacheDir := ".cache"
	app := "bfs"
	file := "token.json"
	want := filepath.Join(home, cacheDir, app, file)
	got := utils.MustGetJWTFile()
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func failOnErr(t *testing.T, e error) {
	t.Helper()
	if e != nil {
		t.Fatalf("failed with error: %q", e)
	}
}

func setConfig(t *testing.T, cfg utils.CLIConfig) {
	t.Helper()

	f := filepath.Join(t.TempDir(), "config.json")

	data, err := json.MarshalIndent(cfg, "", " ")
	failOnErr(t, err)

	failOnErr(t, os.WriteFile(f, data, 0o600))

	utils.ConfigFilePath = f
}

func setJWTToken(t *testing.T, token schemas.JsonToken) {
	t.Helper()

	f := filepath.Join(t.TempDir(), "token.json")

	data, err := json.MarshalIndent(token, "", " ")
	failOnErr(t, err)

	failOnErr(t, os.WriteFile(f, data, 0o600))

	utils.JWTFilePath = f
}

func TestMustReadConfig(t *testing.T) {
	want := utils.CLIConfig{
		URL:      "localhost:8080",
		Username: "ACE",
	}

	setConfig(t, want)

	got := utils.MustReadConfig()

	if want.URL != got.URL {
		t.Fatalf("want %s, got %s", want.URL, got.URL)
	}

	if want.Username != got.Username {
		t.Fatalf("want %s, got %s", want.Username, got.Username)
	}
}

func TestMustReadJWTToken(t *testing.T) {
	want := schemas.JsonToken{
		Token: "super-secret-test-token",
	}

	setJWTToken(t, want)

	got := utils.MustReadJWTToken()

	if want.Token != got.Token {
		t.Fatalf("want %s, got %s", want.Token, got.Token)
	}
}

func TestMustSetJWTToken(t *testing.T) {
	original := schemas.JsonToken{
		Token: "super-secret-test-token-1",
	}

	setJWTToken(t, original)

	got := utils.MustReadJWTToken()

	if original.Token != got.Token {
		t.Fatalf("want %s, got %s", original.Token, got.Token)
	}

	want := schemas.JsonToken{
		Token: "super-secret-test-token-2",
	}

	utils.MustSetJWTToken(want)

	got = utils.MustReadJWTToken()

	if want.Token != got.Token {
		t.Fatalf("want %s, got %s", want.Token, got.Token)
	}
}
