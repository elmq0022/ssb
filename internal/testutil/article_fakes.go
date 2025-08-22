package testutil

import (
	"errors"
	"github.com/google/uuid"
	"ssb/internal/models"
	"ssb/internal/schemas"
)

type FakeArticleRepository struct {
	ArticleStore map[string]models.Article
	UserStore    map[string]models.User
}

func (f *FakeArticleRepository) GetByID(id string) (schemas.ArticleWithAuthorSchema, error) {
	if _, exists := f.ArticleStore[id]; !exists {
		return schemas.ArticleWithAuthorSchema{}, errors.New("Article Not Found")
	}
	a := f.ArticleStore[id]

	if _, exists := f.UserStore[a.Author]; !exists {
		return schemas.ArticleWithAuthorSchema{}, errors.New("User Not Found")
	}

	u := f.UserStore[a.Author]
	b := schemas.UserBrief{
		UserName:  u.UserName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}

	repsonse := schemas.ArticleWithAuthorSchema{
		Title:  a.Title,
		Body:   a.Body,
		Author: b,
	}
	return repsonse, nil
}

func (f *FakeArticleRepository) ListAll() ([]schemas.ArticleWithAuthorSchema, error) {
	var response []schemas.ArticleWithAuthorSchema

	for k, _ := range f.ArticleStore {
		v, err := f.GetByID(k)
		if err != nil {
			return []schemas.ArticleWithAuthorSchema{}, errors.New("User not found for artile")
		}
		response = append(response, v)
	}
	return response, nil
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
	f.ArticleStore[id] = article
	return id, nil
}

func (f *FakeArticleRepository) Update(id string, update schemas.ArticleUpdateSchema) error {
	article, ok := f.ArticleStore[id]
	if !ok {
		return errors.New("article not found")
	}

	if update.Title != nil {
		article.Title = *update.Title
	}

	if update.Body != nil {
		article.Body = *update.Body
	}
	f.ArticleStore[id] = article
	return nil
}

func (f *FakeArticleRepository) Delete(id string) error {
	_, exists := f.ArticleStore[id]
	if exists {
		delete(f.ArticleStore, id)
		return nil
	} else {
		return errors.New("Does not exist")
	}
}

func NewFakeArticleRepository(articles []models.Article, users []models.User) *FakeArticleRepository {
	as := make(map[string]models.Article)
	us := make(map[string]models.User)

	for _, article := range articles {
		as[article.ID] = article
	}

	for _, user := range users {
		us[user.UserName] = user
	}

	f := FakeArticleRepository{
		ArticleStore: as,
		UserStore:    us,
	}

	return &f
}
