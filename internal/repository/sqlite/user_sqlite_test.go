package repo_test

import (
	"database/sql"
	"ssb/internal/auth"
	"ssb/internal/domain/models"
	"ssb/internal/dto"
	"ssb/internal/repository/sqlite"
	"ssb/internal/testutil"
	"testing"
)

func NewUserTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("could not create in memory test db")
	}
	return db
}

func InsertUserIntoDB(t *testing.T,
	db *sql.DB,
	userName,
	firstName,
	lastName,
	email,
	hashedPassword string,
	isActive bool,
	createdAt,
	updatedAt int64,
) {
	t.Helper()
	q := `INSERT INTO users (
      user_name,
      first_name,
      last_name,
      email,
      hashed_password,
	  is_active,
      created_at,
	  updated_at
	) VALUES (
      ?, ?, ?, ?, ?, ?, ?, ?
	)`
	_, err := db.Exec(
		q, userName, firstName, lastName,
		email, hashedPassword, isActive,
		createdAt, updatedAt,
	)
	if err != nil {
		t.Fatalf("could not insert user for test due to error: %v", err)
	}
}

func TestGetUserByUserName(t *testing.T) {
	ur, db := repo.NewUserSqliteRepo(repo.NewTestDB(), testutil.Fc0)
	q := `INSERT
	INTO users (
		user_name,
		first_name,
		last_name,
		email,
		hashed_password,
		is_active,
		created_at,
		updated_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	userName := "tyler.durden"
	firstName := "first name"
	lastName := "last name"
	email := "email@test.com"
	hashedPassword := "random-hashed-value"
	isActive := true
	createdAt := testutil.Fc0.FixedTime.UTC().Unix()
	updatedAt := testutil.Fc0.FixedTime.UTC().Unix()
	_, err := db.Exec(
		q, userName, firstName, lastName,
		email, hashedPassword, isActive,
		createdAt, updatedAt,
	)
	if err != nil {
		t.Fatalf("could not insert user for test due to error: %v", err)
	}

	user, err := ur.GetByUserName(userName)
	if err != nil {
		t.Fatalf("could not get user by user name due to error: %v", err)
	}

	if user.UserName != userName {
		t.Errorf("want %s got %s", userName, user.UserName)
	}
}

func TestCreateUser(t *testing.T) {
	var userName string
	var err error

	data := dto.CreateUserDTO{
		UserName:  "testUser",
		FirstName: "testFirstName",
		LastName:  "testLastName",
		Email:     "test@email.me",
		Password:  "testPassword",
	}

	ur, db := repo.NewUserSqliteRepo(repo.NewTestDB(), testutil.Fc0)
	userName, err = ur.Create(data)

	if err != nil {
		t.Fatalf("%v", err)
	}

	q := `
    SELECT
      user_name,
      first_name,
      last_name,
      email,
      hashed_password,
      created_at,
	  updated_at
    FROM users
    WHERE user_name = ?`

	user := models.User{}
	if err := db.QueryRow(q, userName).Scan(
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		t.Fatalf("%v", err)
	}
	if data.UserName != userName {
		t.Errorf("wanted: %s, got %s", data.UserName, userName)
	}
	if data.FirstName != user.FirstName {
		t.Errorf("wanted: %s, got %s", data.FirstName, user.FirstName)
	}
	if data.LastName != user.LastName {
		t.Errorf("wanted: %s, got %s", data.LastName, user.LastName)
	}
	if data.Email != user.Email {
		t.Errorf("wanted: %s, got %s", data.Email, user.Email)
	}
	ok, err := auth.CheckPassword(data.Password, user.HashedPassword)
	if !ok {
		t.Errorf("password %s did not match hash %s", data.Password, user.HashedPassword)
	}
	if err != nil {
		t.Errorf(
			"password: %s did not match hash %s, got error %v",
			data.Password,
			user.HashedPassword,
			err,
		)
	}
	if user.CreatedAt != testutil.Fc0.FixedTime.Unix() {
		t.Errorf("want %d, got %d", testutil.Fc0.FixedTime.Unix(), user.CreatedAt)
	}
	if user.UpdatedAt != testutil.Fc0.FixedTime.Unix() {
		t.Errorf("want %d, got %d", testutil.Fc0.FixedTime.Unix(), user.UpdatedAt)
	}
}

func TestDeleteUser(t *testing.T) {
	r, db := repo.NewUserSqliteRepo(repo.NewTestDB(), testutil.Fc0)

	userName := "tyler.durden"
	firstName := "first name"
	lastName := "last name"
	email := "email@test.com"
	hashedPassword := "random-hashed-value"
	isActive := true
	createdAt := testutil.Fc0.FixedTime.UTC().Unix()
	updatedAt := testutil.Fc0.FixedTime.UTC().Unix()
	InsertUserIntoDB(
		t,
		db,
		userName,
		firstName,
		lastName,
		email,
		hashedPassword,
		isActive,
		createdAt,
		updatedAt,
	)
	if err := r.Delete(userName); err != nil {
		t.Fatalf("could not delete user due to error: %v", err)
	}

	sql := `SELECT COUNT(*) FROM users WHERE user_name = ?`
	var count int32
	if err := db.QueryRow(sql, userName).Scan(&count); err != nil {
		t.Fatalf("could not query count due to error: %v", err)
	}
	if count != 0 {
		t.Fatalf("wanted 0 users for user name %s, but got %d", userName, count)
	}
}

//TODO: make this a table driven test
func TestUserUpdate(t *testing.T) {
	ur, db := repo.NewUserSqliteRepo(repo.NewTestDB(), testutil.Fc5)

	userName := "tyler.durden"
	firstName := "first name"
	lastName := "last name"
	email := "email@test.com"
	hashedPassword := "random-hashed-value"
	isActive := true
	createdAt := testutil.Fc0.FixedTime.UTC().Unix()
	updatedAt := testutil.Fc0.FixedTime.UTC().Unix()
	InsertUserIntoDB(
		t,
		db,
		userName,
		firstName,
		lastName,
		email,
		hashedPassword,
		isActive,
		createdAt,
		updatedAt,
	)

	newPassword := "new-hashed-password"
	data := dto.UpdateUserDTO{
		UserName:  nil,
		FirstName: nil,
		LastName:  nil,
		Email:     nil,
		Password:  &newPassword,
		IsActive:  nil,
	}
	if err := ur.Update(userName, data); err != nil {
		t.Fatalf("could not updated user due to error: %v", err)
	}
	sql := `SELECT
	  user_name,
	  hashed_password,
	  created_at,
	  updated_at
	FROM users
	WHERE user_name = ?
	`
	user := models.User{}
	if err := db.QueryRow(sql, userName).Scan(
		&user.UserName,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		t.Fatalf("could not get user due to error: %v", err)
	}
	match, err := auth.CheckPassword(newPassword, user.HashedPassword)
	if err != nil {
		t.Fatalf("password check failed due to: %v", err)
	}
	if !match {
		t.Errorf("hash does not match password %s", newPassword)
	}
	if user.CreatedAt != testutil.Fc0.FixedTime.Unix() {
		t.Error("created_at was changed")
	}
	if user.UpdatedAt != testutil.Fc5.FixedTime.Unix() {
		t.Error("updated_at was not updated")
	}
}
