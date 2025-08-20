package repo

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"ssb/internal/models"
	"ssb/internal/schemas"
	"ssb/internal/timeutil"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var defaultArticle = models.Article{}
var defaultError = fmt.Errorf("DefaultError")

func timeFromString(ts string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

type SqliteArticleRepo struct {
	db *sql.DB
	fc timeutil.Clock
}

func NewSqliteArticleRepo(db *sql.DB, clock timeutil.Clock) SqliteArticleRepo {
	return SqliteArticleRepo{
		db: db,
		fc: clock,
	}
}

//go:embed sql/get_article_by_id.sql
var getByIdSQL string

//go:embed sql/list_all_articles.sql
var listAllArticlesSQL string

//go:embed sql/create_article.sql
var createArtcleSQL string

//go:embed sql/delete_article.sql
var deleteArticleSQL string

func (r *SqliteArticleRepo) GetByID(id string) (models.Article, error) {
	var _id string
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

func (r *SqliteArticleRepo) Create(a dto.ArticleCreateDTO) (string, error) {
	id := uuid.New().String()
	now := r.fc.Now().UTC().Format(time.RFC3339Nano)

	_, err := r.db.Exec(createArtcleSQL, id, a.Title, a.Author, a.Body, now, now)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *SqliteArticleRepo) Update(id string, update dto.ArticleUpdateDTO) error {

	var sets []string
	var args []any

	if update.Title != nil {
		sets = append(sets, "title =?")
		args = append(args, update.Title)
	}

	if update.Author != nil {
		sets = append(sets, "author = ?")
		args = append(args, update.Author)
	}

	if update.Body != nil {
		sets = append(sets, "body = ?")
		args = append(args, update.Body)
	}

	sets = append(sets, "updated_at = ?")
	args = append(args, r.fc.Now().Format(time.RFC3339Nano))

	args = append(args, id)

	query := fmt.Sprintf("UPDATE articles SET %s WHERE id = ?", strings.Join(sets, ", "))
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *SqliteArticleRepo) Delete(id string) error {
	_, err := r.db.Exec(deleteArticleSQL, id)
	return err
}
