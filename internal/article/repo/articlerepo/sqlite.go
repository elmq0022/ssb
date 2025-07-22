package repo

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"ssb/internal/api/dto"
	"ssb/internal/article"
)

var defaultArticle = article.Article{}
var defaultError = fmt.Errorf("DefaultError")

type SqliteArticleRepo struct {
	db *sql.DB
}

func GetByID(id string) (article.Article, error) {
	a := article.Article{}
	return a, defaultError
}

func Listall() ([]article.Article, error) {
	articles := [...]article.Article{defaultArticle}
	return articles[:], defaultError
}

func Create(a article.Article) (string, error) {
	return "", defaultError
}

func Update(id string, update dto.ArticleUpdateDTO) error {
	return defaultError
}

func Delete(id string) error {
	return defaultError
}
