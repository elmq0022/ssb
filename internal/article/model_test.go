package article_test

import (
	"ssb/internal/api/dto"
	"ssb/internal/article"
	"ssb/internal/timeutil"
	"testing"
	"time"
)

var now = time.Now()
var later = now.Add(5 * time.Minute)
var fc0 = timeutil.FakeClock{FixedTime: now}
var fc5 = timeutil.FakeClock{FixedTime: later}

const defaultId = 123
const defaultTitle = "defaultTitle"
const defaultAuthor = "defaultAuthor"
const defaultBody = "defaultBody"

var defaultTime = fc0.Now()

func defaultArticle() article.Article {
	return article.NewArticle(
		defaultId, defaultTitle, defaultAuthor, defaultBody, fc0)
}

func TestNewArticle(t *testing.T) {
	want := article.Article{
		Id:          defaultId,
		Title:       defaultTitle,
		Author:      defaultAuthor,
		Body:        defaultBody,
		PublishedAt: defaultTime,
		UpdatedAt:   defaultTime,
	}

	// defaultArticle uses NewArticle to create the article
	got := defaultArticle()
	assertArticleEqual(t, got, want)
}

func TestCloneArticle(t *testing.T) {
	want := defaultArticle()
	got := want.CloneArticle()
	assertArticleEqual(t, got, want)
}

func StringPtr(s string) *string {
	return &s
}

func assertArticleEqual(t *testing.T, got, want article.Article) {
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

func TestUpdateArticleFromDTO(t *testing.T) {
	o := defaultArticle()
	tests := []struct {
		name  string
		dto   dto.ArticleUpdateDTO
		want  article.Article
		clock timeutil.Clock
	}{
		{
			name: "update all fields",
			dto: dto.ArticleUpdateDTO{
				Title:  StringPtr("newTitle"),
				Body:   StringPtr("newBody"),
				Author: StringPtr("newAuthor"),
			},
			want: article.Article{
				Id:          o.Id,
				Title:       "newTitle",
				Body:        "newBody",
				Author:      "newAuthor",
				PublishedAt: now,
				UpdatedAt:   later,
			},
			clock: fc5,
		},
		{
			name: "update author field",
			dto: dto.ArticleUpdateDTO{
				Author: StringPtr("newAuthor"),
			},
			want: article.Article{
				Id:          o.Id,
				Title:       o.Title,
				Body:        o.Body,
				Author:      "newAuthor",
				PublishedAt: now,
				UpdatedAt:   later,
			},
			clock: fc5,
		},
		{
			name: "no updates",
			dto:  dto.ArticleUpdateDTO{},
			want: article.Article{
				Id:          o.Id,
				Title:       o.Title,
				Body:        o.Body,
				Author:      o.Author,
				PublishedAt: now,
				UpdatedAt:   later,
			},
			clock: fc5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := o.CloneArticle()
			got.UpdateArticleFromDTO(tt.dto, tt.clock)
			assertArticleEqual(t, got, tt.want)
		})
	}
}
