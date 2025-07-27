package repo

import (
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"ssb/internal/domain/models"
	"ssb/internal/dto"
	"time"
)

var defaultArticle = models.Article{}
var defaultError = fmt.Errorf("DefaultError")

func timeFromString(ts string) time.Time {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

type SqliteArticleRepo struct {
	db *sql.DB
}

func NewSqliteArticleRepo(db *sql.DB) SqliteArticleRepo {
	return SqliteArticleRepo{db: db}
}

//go:embed sql/get_article_by_id.sql
var getByIdSQL string

//go:embed sql/list_all_articles.sql
var listAllArticlesSQL string

func (r *SqliteArticleRepo) GetByID(id uint32) (models.Article, error) {
	var _id uint32
	var title string
	var author string
	var body string
	var publishedAt string
	var updatedAt string

	row := r.db.QueryRow(getByIdSQL, id)
	err := row.Scan(&_id, &title, &author, &body, &publishedAt, &updatedAt)
	if err != nil {
		return models.Article{}, err
	}

	a := models.Article{}
	a.ID = _id
	a.Title = title
	a.Author = author
	a.Body = body
	a.PublishedAt = timeFromString(publishedAt)
	a.UpdatedAt = timeFromString(updatedAt)
	return a, nil
}

func (r *SqliteArticleRepo) ListAll() ([]models.Article, error) {
	rows, err := r.db.Query(listAllArticlesSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article

	for rows.Next() {
		var a models.Article
		var publishedAt string
		var updatedAt string
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Author,
			&a.Body,
			&publishedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		a.PublishedAt = timeFromString(publishedAt)
		a.UpdatedAt = timeFromString(updatedAt)
		articles = append(articles, a)
	}

	return articles, nil
}

func (r *SqliteArticleRepo) Create(a models.Article) (string, error) {
	return "", defaultError
}

func (r *SqliteArticleRepo) Update(id string, update dto.ArticleUpdateDTO) error {
	return defaultError
}

func (r *SqliteArticleRepo) Delete(id string) error {
	return defaultError
}
