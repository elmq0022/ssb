package repo_test

import (
	"database/sql"
	tdb "ssb/internal/db"
	"ssb/internal/models"
	"ssb/internal/repo/sqlite"
	"ssb/internal/schemas"
	"ssb/internal/testutil"
	"ssb/internal/timeutil"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// TODO move to test utils
func asserEqual(t *testing.T, want, got any) {
	t.Helper()
	if !cmp.Equal(want, got) {
		t.Errorf("mismatch (-want +got):\n%s", cmp.Diff(want, got))
	}
}

const INSERT_ARTICLE = `
INSERT INTO articles (
	id, 
	title,
	author,
	body,
	published_at,
	updated_at
) VALUES (?, ?, ?, ?, ?, ?)`

const INSERT_USER = `
INSERT INTO users (
  user_name,
  first_name,
  last_name,
  email,
  hashed_password,
  is_active,
  created_at,
  updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

func NewTestRepo(clock timeutil.Clock) (repo.SqliteArticleRepo, *sql.DB) {
	db := tdb.MustNewTestDB()
	r := repo.NewSqliteArticleRepo(db, clock)
	return r, db
}

func TestNewDBReturnsZeroRows(t *testing.T) {
	db := tdb.MustNewTestDB()
	defer db.Close()

	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM articles").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatalf("wanted: 0 rows but got: %d rows", count)
	}
}

func mustCreateUser(t *testing.T, db *sql.DB, userName string) {
	t.Helper()
	_, err := db.Exec(
		INSERT_USER,
		userName,
		"first_name",
		"last_name",
		"test@example.com",
		"super_secret",
		true,
		testutil.Fc0.FixedTime.Unix(),
		testutil.Fc0.FixedTime.Unix(),
	)

	if err != nil {
		t.Fatalf("could not create user: %v", err)
	}
}

func TestGetArticleByID(t *testing.T) {
	r, db := NewTestRepo(testutil.Fc0)
	defer db.Close()

	want := testutil.NewArticle(testutil.Fc0)
	mustCreateUser(t, db, want.Author)
	db.Exec(INSERT_ARTICLE,
		want.ID,
		want.Title,
		want.Author,
		want.Body,
		want.PublishedAt.UTC().Format(time.RFC3339Nano),
		want.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)

	got, err := r.GetByID(want.ID)

	if err != nil {
		t.Fatal(err)
	}

	testutil.AssertArticleEqual(t, got, want)
}

func TestGetAllArticles(t *testing.T) {
	r, db := NewTestRepo(testutil.Fc0)
	defer db.Close()

	a1 := testutil.NewArticle(testutil.Fc0)

	a2 := testutil.NewArticle(
		testutil.Fc5,
		testutil.WithID("2"),
		testutil.WithAuthor("Author 2"),
		testutil.WithTitle("Title 2"),
		testutil.WithBody("Body 2"),
	)

	want := []models.Article{a1, a2}
	for _, a := range want {
		db.Exec(INSERT_ARTICLE,
			a.ID,
			a.Title,
			a.Author,
			a.Body,
			a.PublishedAt.UTC().Format(time.RFC3339Nano),
			a.UpdatedAt.UTC().Format(time.RFC3339Nano),
		)
	}

	got, err := r.ListAll()
	if err != nil {
		t.Fatalf("%v", err)
	}

	asserEqual(t, want, got)
}

func TestCreateArticle(t *testing.T) {
	r, db := NewTestRepo(testutil.Fc0)
	defer db.Close()

	title := "New Title"
	author := "New Author"
	body := "New Body"
	create_dto := schemas.ArticleCreateSchema{
		Title:    title,
		UserName: author,
		Body:     body,
	}

	id, err := r.Create(create_dto)
	if err != nil {
		t.Fatalf("%q", err)
	}

	want := testutil.NewArticle(
		testutil.Fc0,
		testutil.WithID(id),
		testutil.WithTitle(title),
		testutil.WithAuthor(author),
		testutil.WithBody(body),
	)

	// TODO: I should probably use raw sql here.
	got, err := r.GetByID(id)
	asserEqual(t, want, got)
}

func ptrFromString(s string) *string {
	return &s
}

// TODO: Need to populated DB with an article.
// Then do the update and then check the result.

/*
func TestUpdateArticle(t *testing.T) {
	tests := []struct {
		name    string
		a       models.Article
		updates schemas.ArticleUpdateSchema
		want    models.Article
	}{
		{
			"no-op",
			testutil.NewArticle(testutil.Fc0),
			schemas.ArticleUpdateSchema{Title: nil, UserName: nil, Body: nil},
			testutil.NewArticle(
				testutil.Fc0,
				testutil.WithUpdatedAt(testutil.Fc5)),
		},
		{
			"update-all",
			testutil.NewArticle(testutil.Fc0),
			schemas.ArticleUpdateSchema{
				Title:    ptrFromString("newTitle"),
				UserName: ptrFromString("newAuthor"),
				Body:     ptrFromString("newBody"),
			},
			testutil.NewArticle(
				testutil.Fc0,
				testutil.WithTitle("newTitle"),
				testutil.WithAuthor("newAuthor"),
				testutil.WithBody("newBody"),
				testutil.WithUpdatedAt(testutil.Fc5)),
		},
		{
			"update-title",
			testutil.NewArticle(testutil.Fc0),
			schemas.ArticleUpdateSchema{
				Title:    ptrFromString("newTitle"),
				UserName: nil,
				Body:     nil,
			},
			testutil.NewArticle(
				testutil.Fc0,
				testutil.WithTitle("newTitle"),
				testutil.WithUpdatedAt(testutil.Fc5),
			),
		},
		{
			"update-author",
			testutil.NewArticle(testutil.Fc0),
			schemas.ArticleUpdateSchema{
				Title:    nil,
				UserName: ptrFromString("newAuthor"),
				Body:     nil,
			},
			testutil.NewArticle(
				testutil.Fc0,
				testutil.WithAuthor("newAuthor"),
				testutil.WithUpdatedAt(testutil.Fc5),
			),
		},
		{
			"update-body",
			testutil.NewArticle(testutil.Fc0),
			schemas.ArticleUpdateSchema{
				Title:    nil,
				UserName: nil,
				Body:     ptrFromString("newBody"),
			},
			testutil.NewArticle(
				testutil.Fc0,
				testutil.WithBody("newBody"),
				testutil.WithUpdatedAt(testutil.Fc5),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, db := NewTestRepo(testutil.Fc5)

			t.Cleanup(func() { db.Close() })

			sql := `INSERT INTO articles (
				id,
				title,
				author,
				body,
				published_at,
				updated_at
			)
			VALUES (?, ?, ?, ?, ?, ?)`

			_, err := db.Exec(sql,
				tt.a.ID,
				tt.a.Title,
				tt.a.Author,
				tt.a.Body,
				tt.a.PublishedAt.Format(time.RFC3339Nano),
				tt.a.UpdatedAt.Format(time.RFC3339Nano),
			)

			if err != nil {
				t.Fatalf("Could not create update target in db. Error: %q", err)
			}

			r.Update(tt.a.ID, tt.updates)
			got, err := r.GetByID(tt.a.ID)
			if err != nil {
				t.Fatalf("%q", err)
			}

			asserEqual(t, tt.want, got)
		})
	}
}

func TestDeleteArticle(t *testing.T) {
	r, db := NewTestRepo(testutil.Fc0)
	t.Cleanup(func() { db.Close() })
	id := "delete-me"
	a := testutil.NewArticle(testutil.Fc0, testutil.WithID(id))
	sql := `INSERT INTO articles (
		id,
		title,
		author,
		body,
		published_at,
		updated_at
	)
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(
		sql,
		a.ID,
		a.Title,
		a.Author,
		a.Body,
		a.PublishedAt,
		a.UpdatedAt,
	)

	if err != nil {
		t.Fatalf("%q", err)
	}

	err = r.Delete(id)
	if err != nil {
		t.Fatalf("%q", err)
	}

	want := 0
	check := `SELECT COUNT(*) FROM articles WHERE id = ?`

	var got int
	err = db.QueryRow(check, a.ID).Scan(&got)
	if err != nil {
		t.Fatalf("%q", err)
	}
	if want != got {
		t.Errorf("want: %d, got %d", want, got)
	}
}
*/
