package article_test

import (
	"ssb/internal/article"
	"ssb/internal/timeutil"
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

func TestNewArticleFunc(t *testing.T) {
	fixedTime := time.Now()
	fc := timeutil.FakeClock{fixedTime}

	want := article.Article{
		Id:          123,
		Title:       "title",
		Author:      "author",
		Body:        "body",
		PublishedAt: fixedTime,
		UpdatedAt:   fixedTime,
	}

	got := article.NewArticle(123, "title", "author", "body", fc)

	if want != got {
		t.Fatalf("want: %v but got: %v", want, got)
	}
}
