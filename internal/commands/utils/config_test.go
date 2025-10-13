package utils_test

import (
	"os"
	"path/filepath"
	tu "ssb/internal/commands/testUtils"
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

func TestMustReadConfig(t *testing.T) {
	want := utils.CLIConfig{
		URL:      "localhost:8080",
		Username: "ACE",
	}

	tu.SetConfig(t, want)

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

	tu.SetJWTToken(t, want)

	got := utils.MustReadJWTToken()

	if want.Token != got.Token {
		t.Fatalf("want %s, got %s", want.Token, got.Token)
	}
}

func TestMustSetJWTToken(t *testing.T) {
	original := schemas.JsonToken{
		Token: "super-secret-test-token-1",
	}

	tu.SetJWTToken(t, original)

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

func TestMustSetConfig(t *testing.T) {
	orignal := utils.CLIConfig{
		URL:      "url-1",
		Username: "user-1",
	}
	tu.SetConfig(t, orignal)
	got := utils.MustReadConfig()

	if orignal.URL != got.URL || orignal.Username != got.Username {
		t.Fatalf("want %q, got %q", orignal, got)
	}

	want := utils.CLIConfig{
		URL:      "url-2",
		Username: "user-2",
	}

	utils.MustSetConfig(&want)
	got = utils.MustReadConfig()

	if want.URL != got.URL || want.Username != got.Username {
		t.Fatalf("want %q, got %q", want, got)
	}
}
