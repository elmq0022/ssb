package repo

import (
	"ssb/internal/models"
	"ssb/internal/schemas"
)

type ArticleRepository interface {
	GetByID(id string) (schemas.ArticleWithAuthorSchema, error)
	ListAll() ([]schemas.ArticleWithAuthorSchema, error)
	Create(a schemas.ArticleCreateSchema) (string, error)
	Update(id string, update schemas.ArticleUpdateSchema) error
	Delete(id string) error
}

type UserRepository interface {
	GetByUserName(username string) (models.User, error)
	Create(data schemas.CreateUserDTO) (string, error)
	Update(userName string, data schemas.UpdateUserDTO) error
	Delete(userName string) error
}
