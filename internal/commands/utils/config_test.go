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
