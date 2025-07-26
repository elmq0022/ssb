package testutil_test

import (
	"ssb/internal/article"
	"ssb/internal/testutil"
	"testing"
)

func TestNewArticle(t *testing.T) {
	want := article.Article{
		Id:          testutil.DefaultId,
		Title:       testutil.DefaultTitle,
		Author:      testutil.DefaultAuthor,
		Body:        testutil.DefaultBody,
		PublishedAt: testutil.DefaultTime,
		UpdatedAt:   testutil.DefaultTime,
	}

	// defaultArticle uses NewArticle to create the article
	got := testutil.DefaultArticle()
	testutil.AssertArticleEqual(t, got, want)
}
