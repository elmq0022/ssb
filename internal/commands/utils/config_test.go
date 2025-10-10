package utils_test

import (
	"os"
	"path/filepath"
	"ssb/internal/commands/utils"
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
