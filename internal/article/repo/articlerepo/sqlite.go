package repo

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"ssb/internal/api/dto"
	"ssb/internal/article"
	"time"
)

var defaultArticle = article.Article{}
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

func (r *SqliteArticleRepo) GetByID(id int32) (article.Article, error) {
	var _id int32
	var title string
	var author string
	var body string
	var publishedAt string
	var updatedAt string

	sql := `SELECT 
		id, 
		title,
		author, 
		body, 
		published_at, 
		updated_at 
	FROM ARTICLES 
	WHERE id = ?`

	row := r.db.QueryRow(sql, id)
	err := row.Scan(&_id, &title, &author, &body, &publishedAt, &updatedAt)
	if err != nil {
		return article.Article{}, err
	}

	a := article.Article{}
	a.Id = _id
	a.Title = title
	a.Author = author
	a.Body = body
	a.PublishedAt = timeFromString(publishedAt)
	a.UpdatedAt = timeFromString(updatedAt)
	return a, nil
}

func (r *SqliteArticleRepo) Listall() ([]article.Article, error) {
	articles := [...]article.Article{defaultArticle}
	return articles[:], defaultError
}

func (r *SqliteArticleRepo) Create(a article.Article) (string, error) {
	return "", defaultError
}

func (r *SqliteArticleRepo) Update(id string, update dto.ArticleUpdateDTO) error {
	return defaultError
}

func (r *SqliteArticleRepo) Delete(id string) error {
	return defaultError
}
