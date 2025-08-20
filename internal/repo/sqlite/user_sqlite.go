package repo

import (
	"database/sql"
	"log"
	"ssb/internal/pkg/auth"
	"ssb/internal/models"
	"ssb/internal/schemas"
	"ssb/internal/timeutil"

	sq "github.com/Masterminds/squirrel"
)

var schema string = `
CREATE TABLE users (
	pk INTEGER PRIMARY KEY AUTOINCREMENT,
	user_name TEXT UNIQUE NOT NULL,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	email TEXT UNIQUE NOT NULL,
	hashed_password TEXT NOT NULL,
	is_active BOOL NOT NULL,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL
)`

type UserSqliteRepo struct {
	db    *sql.DB
	clock timeutil.Clock
}

func NewTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("failed to create schema %v", err)
	}
	return db
}

func NewUserSqliteRepo(db *sql.DB, clock timeutil.Clock) (UserSqliteRepo, *sql.DB) {
	r := UserSqliteRepo{
		db:    db,
		clock: clock,
	}
	return r, db
}

func (r *UserSqliteRepo) GetByUserName(userName string) (models.User, error) {
	sql := sq.Select(
		"user_name",
		"first_name",
		"last_name",
		"email",
		"hashed_password",
		"is_active",
		"created_at",
		"updated_at",
	).From("users").Where(sq.Eq{"user_name": userName})
	row := sql.RunWith(r.db).QueryRow()
	user := models.User{}
	err := row.Scan(
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserSqliteRepo) Create(data dto.CreateUserDTO) (string, error) {
	now := r.clock.Now().UTC().Unix()
	hashedPassword := auth.HashPassword(data.Password)
	sql := sq.Insert("users").Columns(
		"user_name",
		"first_name",
		"last_name",
		"email",
		"hashed_password",
		"is_active",
		"created_at",
		"updated_at",
	).Values(
		data.UserName,
		data.FirstName,
		data.LastName,
		data.Email,
		hashedPassword,
		true,
		now,
		now,
	)
	_, err := sql.RunWith(r.db).Exec()
	if err != nil {
		return "", err
	}
	return data.UserName, nil
}

func (r *UserSqliteRepo) Update(userName string, data dto.UpdateUserDTO) error {
	q := sq.Update("users")
	if data.UserName != nil {
		q = q.Set("user_name", *data.UserName)
	}
	if data.FirstName != nil {
		q = q.Set("first_name", *data.FirstName)
	}
	if data.LastName != nil {
		q = q.Set("last_name", *data.LastName)
	}
	if data.Email != nil {
		q = q.Set("email", *data.Email)
	}
	if data.IsActive != nil {
		q = q.Set("is_active", *data.IsActive)
	}
	if data.Password != nil {
		q = q.Set("hashed_password", auth.HashPassword(*data.Password))
	}
	q = q.Set("updated_at", r.clock.Now().UTC().Unix())
	q = q.Where(sq.Eq{"user_name": userName})
	_, err := q.RunWith(r.db).Exec()
	return err
}

func (r *UserSqliteRepo) Delete(userName string) error {
	sql := sq.Delete("users").Where(sq.Eq{"user_name": userName})
	_, err := sql.RunWith(r.db).Exec()
	return err
}
