package repo

import (
	"ssb/internal/api/dto"
	"ssb/internal/article"
)

type ArticleRespository interface {
	GetByID(id string) (article.Article, error)
	Listall() ([]article.Article, error)
	Create(a article.Article) (string, error)
	Update(id string, update dto.ArticleUpdateDTO) error
	Delete(id string) error
}
