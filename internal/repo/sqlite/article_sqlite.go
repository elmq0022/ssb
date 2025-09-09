package repo

import (
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"ssb/internal/models"
	"ssb/internal/schemas"
	"ssb/internal/timeutil"
	"strings"
	"time"
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
	db *sqlx.DB
	fc timeutil.Clock
}

func NewSqliteArticleRepo(db *sqlx.DB, clock timeutil.Clock) SqliteArticleRepo {
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

func (r *SqliteArticleRepo) GetByID(id string) (schemas.ArticleWithAuthorSchema, error) {
	a := schemas.ArticleWithAuthorSchema{}
	if err := r.db.Get(&a, getByIdSQL, id); err != nil {
		return schemas.ArticleWithAuthorSchema{}, err
	}
	return a, nil
}

func (r *SqliteArticleRepo) ListAll() ([]schemas.ArticleWithAuthorSchema, error) {
	var a []schemas.ArticleWithAuthorSchema
	if err := r.db.Select(&a, listAllArticlesSQL); err != nil {
		return nil, err
	}
	return a, nil
}

func (r *SqliteArticleRepo) Create(username string, a schemas.ArticleCreateSchema) (string, error) {
	id := uuid.New().String()
	now := r.fc.Now().UTC().Format(time.RFC3339Nano)

	_, err := r.db.Exec(createArtcleSQL, id, a.Title, username, a.Body, now, now)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *SqliteArticleRepo) Update(id string, update schemas.ArticleUpdateSchema) error {

	var sets []string
	var args []any

	if update.Title != nil {
		sets = append(sets, "title =?")
		args = append(args, update.Title)
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
