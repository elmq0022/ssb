package repo

import (
	"ssb/internal/domain/models"
	"ssb/internal/dto"
)

type ArticleRespository interface {
	GetByID(id string) (models.Article, error)
	ListAll() ([]models.Article, error)
	Create(a models.Article) (string, error)
	Update(id string, update dto.ArticleUpdateDTO) error
	Delete(id string) error
}
