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

const DefaultId = 123
const DefaultTitle = "defaultTitle"
const DefaultAuthor = "defaultAuthor"
const DefaultBody = "defaultBody"

var DefaultTime = Fc0.Now()

func DefaultArticle() models.Article {
	return models.Article{
		ID:          DefaultId,
		Title:       DefaultTitle,
		Author:      DefaultAuthor,
		Body:        DefaultBody,
		PublishedAt: Now,
		UpdatedAt:   Now,
	}
}
