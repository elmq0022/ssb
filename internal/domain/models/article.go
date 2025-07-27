package models

import (
	"time"
)

type Article struct {
	ID          uint32
	Title       string
	Author      string
	Body        string
	PublishedAt time.Time
	UpdatedAt   time.Time
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
