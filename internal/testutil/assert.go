package testutil

import (
	"ssb/internal/article"
	"testing"
)

func AssertArticleEqual(t *testing.T, got, want article.Article) {
	t.Helper()

	if got.Id != want.Id {
		t.Errorf("Id mismatch: got %v, want %v", got.Id, want.Id)
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
