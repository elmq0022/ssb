package utils_test

import (
	"bytes"
	"ssb/internal/commands/utils"
	"strings"
	"testing"
)

func TestPrompt(t *testing.T) {
	want := "test input"
	r := strings.NewReader(want)
	var w bytes.Buffer
	msg := "test prompt: "
	got := utils.Prompt(r, &w, msg)
	if want != got {
		t.Fatalf("want %s, got %s", want, got)
	}
}

func TestReadPasswordTwice_Success(t *testing.T) {
	i := 0
	inputs := []string{"pwd", "pwd"}
	fakePrompt := func(msg string) string {
		v := inputs[i]
		i++
		return v
	}

	got, err := utils.ReadPasswordTwice(fakePrompt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "pwd"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestReadPasswordTwice_Mismatch(t *testing.T) {
	i := 0
	inputs := []string{"abc", "xyz"}
	fakePrompt := func(msg string) string {
		v := inputs[i]
		i++
		return v
	}

	_, err := utils.ReadPasswordTwice(fakePrompt)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
