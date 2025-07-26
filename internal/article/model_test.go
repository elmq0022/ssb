package article_test

import (
	"ssb/internal/api/dto"
	"ssb/internal/article"
	"ssb/internal/testutil"
	"ssb/internal/timeutil"
	"testing"
)

/*
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
	got := testutil.DefaultArticle()
	testutil.AssertArticleEqual(t, got, want)
}
*/

func TestCloneArticle(t *testing.T) {
	want := testutil.DefaultArticle()
	got := want.CloneArticle()
	testutil.AssertArticleEqual(t, got, want)
}

func StringPtr(s string) *string {
	return &s
}

func TestUpdateArticleFromDTO(t *testing.T) {
	o := testutil.DefaultArticle()
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
				PublishedAt: testutil.Now,
				UpdatedAt:   testutil.Later,
			},
			clock: testutil.Fc5,
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
				PublishedAt: testutil.Now,
				UpdatedAt:   testutil.Later,
			},
			clock: testutil.Fc5,
		},
		{
			name: "no updates",
			dto:  dto.ArticleUpdateDTO{},
			want: article.Article{
				Id:          o.Id,
				Title:       o.Title,
				Body:        o.Body,
				Author:      o.Author,
				PublishedAt: testutil.Now,
				UpdatedAt:   testutil.Later,
			},
			clock: testutil.Fc5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := o.CloneArticle()
			got.UpdateArticleFromDTO(tt.dto, tt.clock)
			testutil.AssertArticleEqual(t, got, tt.want)
		})
	}
}
