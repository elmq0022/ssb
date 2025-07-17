package article

import (
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

func NewArticle(id int32, title, author, body string) Article {
	a := Article{
		Id:          id,
		Title:       title,
		Author:      author,
		Body:        body,
		PublishedAt: time.Now(),
		UpdatedAt:   time.Now(),
	}
	return a
}
