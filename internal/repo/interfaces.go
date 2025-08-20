package repo

import (
	"ssb/internal/models"
	"ssb/internal/schemas"
)

type ArticleRepository interface {
	GetByID(id string) (models.Article, error)
	ListAll() ([]models.Article, error)
	Create(a dto.ArticleCreateDTO) (string, error)
	Update(id string, update dto.ArticleUpdateDTO) error
	Delete(id string) error
}

type UserRepository interface {
	GetByUserName(username string) (models.User, error)
	Create(data dto.CreateUserDTO) (string, error)
	Update(userName string, data dto.UpdateUserDTO) error
	Delete(userName string) error
}
