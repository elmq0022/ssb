package testutil

import (
	"errors"
	"github.com/google/uuid"
	"ssb/internal/models"
	"ssb/internal/schemas"
)

type FakeArticleRepository struct {
	Store map[string]models.Article
}

func (f *FakeArticleRepository) GetByID(id string) (models.Article, error) {
	article, exists := f.Store[id]
	if !exists {
		return models.Article{}, errors.New("Article Not Found")
	} else {
		return article, nil
	}
}

func (f *FakeArticleRepository) ListAll() ([]models.Article, error) {
	var articles []models.Article

	for _, v := range f.Store {
		articles = append(articles, v)
	}
	return articles, nil
}

func (f *FakeArticleRepository) Create(a schemas.ArticleCreateSchema) (string, error) {
	id := uuid.New().String()
	article := NewArticle(
		Fc0,
		WithID(id),
		WithAuthor(a.UserName),
		WithTitle(a.Title),
		WithBody(a.Body),
	)
	f.Store[id] = article
	return id, nil
}

func (f *FakeArticleRepository) Update(id string, update schemas.ArticleUpdateSchema) error {
	article, ok := f.Store[id]
	if !ok {
		return errors.New("article not found")
	}

	if update.Title != nil {
		article.Title = *update.Title
	}

	if update.UserName != nil {
		article.Author = *update.UserName
	}

	if update.Body != nil {
		article.Body = *update.Body
	}
	f.Store[id] = article
	return nil
}

func (f *FakeArticleRepository) Delete(id string) error {
	_, exists := f.Store[id]
	if exists {
		delete(f.Store, id)
		return nil
	} else {
		return errors.New("Does not exist")
	}
}

func NewFakeArticleRepository(articles []models.Article) *FakeArticleRepository {
	s := make(map[string]models.Article)

	for _, article := range articles {
		s[article.ID] = article
	}

	f := FakeArticleRepository{
		Store: s,
	}

	return &f
}
