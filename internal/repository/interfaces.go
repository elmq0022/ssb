package repo

import (
	"ssb/internal/domain/models"
	"ssb/internal/dto"
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
	Update(data dto.UpdateUserDTO) error
	Delete(id string) error
}
