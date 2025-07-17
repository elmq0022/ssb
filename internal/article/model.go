package article

import (
	"ssb/internal/timeutil"
	"time"
)

type Article struct {
	Id          int32
	Title       string
	Author      string
	Body        string
	PublishedAt time.Time
	UpdatedAt   time.Time
}

func NewArticle(id int32, title, author, body string, clock timeutil.Clock) Article {
	a := Article{
		Id:          id,
		Title:       title,
		Author:      author,
		Body:        body,
		PublishedAt: clock.Now(),
		UpdatedAt:   clock.Now(),
	}
	return a
}
