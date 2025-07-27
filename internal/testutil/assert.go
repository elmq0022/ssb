package testutil

import (
	"ssb/internal/domain/models"
	"testing"
)

func AssertArticleEqual(t *testing.T, got, want models.Article) {
	t.Helper()

	if got.ID != want.ID {
		t.Errorf("Id mismatch: got %v, want %v", got.ID, want.ID)
	}
	if got.Title != want.Title {
		t.Errorf("Title mismatch: got %q, want %q", got.Title, want.Title)
	}
	if got.Body != want.Body {
		t.Errorf("Body mismatch: got %q, want %q", got.Body, want.Body)
	}
	if got.Author != want.Author {
		t.Errorf("Author mismatch: got %q, want %q", got.Author, want.Author)
	}
	if !got.PublishedAt.Equal(want.PublishedAt) {
		t.Errorf("PublishedAt mismatch: got %v, want %v", got.PublishedAt, want.PublishedAt)
	}
	if !got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: got %v, want %v", got.UpdatedAt, want.UpdatedAt)
	}
}
