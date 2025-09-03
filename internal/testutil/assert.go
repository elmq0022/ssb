package testutil

import (
	"github.com/google/go-cmp/cmp"
	"ssb/internal/models"
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

// TODO move to test utils
func asserEqual(t *testing.T, want, got any) {
	t.Helper()
	if !cmp.Equal(want, got) {
		t.Errorf("mismatch (-want +got):\n%s", cmp.Diff(want, got))
	}
}
