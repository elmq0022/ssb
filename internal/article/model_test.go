package article_test

import (
	"ssb/internal/article"
	"testing"
	"time"
)

func TestNewArticle(t *testing.T) {
	publishedAt, err := time.Parse(time.RFC3339, "2025-07-14T21:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse %v", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, "2025-07-14T21:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse %v", err)
	}

	a := article.Article{
		Id:          123,
		Title:       "title",
		Author:      "author",
		Body:        "body",
		PublishedAt: publishedAt,
		UpdatedAt:   updatedAt,
	}

	t.Logf("%q", a)
}
