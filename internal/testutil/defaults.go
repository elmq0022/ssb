package testutil

import (
	"ssb/internal/domain/models"
	"ssb/internal/timeutil"
	"time"
)

var Now = time.Now().UTC()
var Later = Now.Add(5 * time.Minute)
var Fc0 = timeutil.FakeClock{FixedTime: Now}
var Fc5 = timeutil.FakeClock{FixedTime: Later}

const DefaultId = "1"
const DefaultTitle = "defaultTitle"
const DefaultAuthor = "defaultAuthor"
const DefaultBody = "defaultBody"

var DefaultTime = Fc0.Now()

type ArticleOpt func(*models.Article)

func NewArticle(clock timeutil.Clock, opts ...ArticleOpt) models.Article {
	now := clock.Now().UTC()
	a := models.Article{
		ID:          DefaultId,
		Title:       DefaultTitle,
		Author:      DefaultAuthor,
		Body:        DefaultBody,
		PublishedAt: now,
		UpdatedAt:   now,
	}

	for _, opt := range opts {
		opt(&a)
	}

	return a
}

func WithID(id string) ArticleOpt {
	return func(a *models.Article) {
		a.ID = id
	}
}

func WithTitle(title string) ArticleOpt {
	return func(a *models.Article) {
		a.Title = title
	}
}

func WithAuthor(author string) ArticleOpt {
	return func(a *models.Article) {
		a.Author = author
	}
}

func WithBody(body string) ArticleOpt {
	return func(a *models.Article) {
		a.Body = body
	}
}

func WithPublishedAt(clock timeutil.Clock) ArticleOpt {
	return func(a *models.Article) {
		a.PublishedAt = clock.Now().UTC()
	}
}

func WithUpdatedAt(clock timeutil.Clock) ArticleOpt {
	return func(a *models.Article) {
		a.UpdatedAt = clock.Now().UTC()
	}
}
