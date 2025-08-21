package models_test

import (
	"ssb/internal/testutil"
	"testing"
)

func TestCloneArticle(t *testing.T) {
	want := testutil.NewArticle(testutil.Fc0)
	got := want.CloneArticle()
	testutil.AssertArticleEqual(t, got, want)
}

func StringPtr(s string) *string {
	return &s
}
