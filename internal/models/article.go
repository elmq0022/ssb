package models

import (
	"time"
)

type Article struct {
	ID          string
	Title       string
	Author      string
	Body        string
	PublishedAt time.Time
	UpdatedAt   time.Time
}

type ArticleWithAuthor struct {
	ID string
	Title string
	Author struct {
		UserName string
		FirstName string
		LastName string
	}
	Body string
	PublishedAt int64
	UpdatedAt int64
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
