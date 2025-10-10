package utils_test

import (
	"os"
	"path/filepath"
	"ssb/internal/commands/utils"
	"testing"
)

func TestMustGetConfigDir(t *testing.T) {
	home, _ := os.UserHomeDir()
	configDir := ".config"
	app := "bfs"
	file := "config.json"
	want := filepath.Join(home, configDir, app, file)
	got := utils.MustGetConfigDir()
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestMustGetCacheDir(t *testing.T) {
	home, _ := os.UserHomeDir()
	cacheDir := ".cache"
	app := "bfs"
	file := "token.json"
	want := filepath.Join(home, cacheDir, app, file)
	got := utils.MustGetCacheDir()
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}
