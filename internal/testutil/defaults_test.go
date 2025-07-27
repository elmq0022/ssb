package testutil_test

import (
	"ssb/internal/domain/models"
	"ssb/internal/testutil"
	"testing"
)

func TestDefaultArticle(t *testing.T) {
	want := models.Article{
		ID:          testutil.DefaultId,
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
