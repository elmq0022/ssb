package models

import (
	"time"
)

type Article struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Author      string    `db:"author"`
	Body        string    `db:"body"`
	PublishedAt time.Time `db:"published_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (a Article) CloneArticle() Article {
	return Article{
		ID:          a.ID,
		Title:       a.Title,
		Author:      a.Author,
		Body:        a.Body,
		PublishedAt: a.PublishedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}
