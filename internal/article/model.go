package article

import (
	"ssb/internal/api/dto"
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
	now := clock.Now()
	a := Article{
		Id:          id,
		Title:       title,
		Author:      author,
		Body:        body,
		PublishedAt: now,
		UpdatedAt:   now,
	}
	return a
}

func (a Article) CloneArticle() Article {
	return Article{
		Id:          a.Id,
		Title:       a.Title,
		Author:      a.Author,
		Body:        a.Body,
		PublishedAt: a.PublishedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func (a *Article) UpdateArticleFromDTO(ad dto.ArticleUpdateDTO, clock timeutil.Clock) {
	if ad.Author != nil {
		a.Author = *ad.Author
	}
	if ad.Title != nil {
		a.Title = *ad.Title
	}
	if ad.Body != nil {
		a.Body = *ad.Body
	}
	a.UpdatedAt = clock.Now()
}
