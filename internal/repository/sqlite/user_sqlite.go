package repo

import (
	"database/sql"
	"errors"
	"log"
	"ssb/internal/domain/models"
	"ssb/internal/dto"
	"ssb/internal/timeutil"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var schema string = `
CREATE TABLE users (
	pk INTEGER PRIMARY KEY AUTOINCREMENT,
	id TEXT UNIQUE NOT NULL,
	user_name TEXT UNIQUE NOT NULL,
	first_name TEXT UNIQUE NOT NULL,
	last_name TEXT UNIQUE NOT NULL,
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

func (r *UserSqliteRepo) GetByID(id string) (models.User, error) {
	return models.User{}, errors.New("Not Implemented")
}

func (r *UserSqliteRepo) Create(data dto.CreateUserDTO) (string, error) {
	now := r.clock.Now().UTC().Unix()
	id := uuid.New().String()

	sql := sq.Insert("users").Columns(
		"id",
		"user_name",
		"first_name",
		"last_name",
		"email",
		"hashed_password",
		"is_active",
		"created_at",
		"updated_at",
	).Values(
		id,
		data.UserName,
		data.FirstName,
		data.LastName,
		data.Email,
		"password",
		true,
		now,
		now,
	)

	_, err := sql.RunWith(r.db).Exec()
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *UserSqliteRepo) Update(data dto.UpdateUserDTO) error {
	return errors.New("Not Implemented")
}

func (r *UserSqliteRepo) Delete(id string) error {
	return errors.New("Not Implemented")
}
