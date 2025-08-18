package repo_test

import (
	"database/sql"
	"ssb/internal/auth"
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
	id,
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
	  id,
      user_name,
      first_name,
      last_name,
      email,
      hashed_password,
	  is_active,
      created_at,
	  updated_at
	) VALUES (
      ?, ?, ?, ?, ?, ?, ?, ?, ?
	)`
	_, err := db.Exec(
		q, id, userName, firstName, lastName,
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
		id,
		user_name,
		first_name,
		last_name,
		email,
		hashed_password,
		is_active,
		created_at,
		updated_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	id := "id"
	userName := "tyler.durden"
	firstName := "first name"
	lastName := "last name"
	email := "email@test.com"
	hashedPassword := "random-hashed-value"
	isActive := true
	createdAt := testutil.Fc0.FixedTime.UTC().Unix()
	updatedAt := testutil.Fc0.FixedTime.UTC().Unix()
	_, err := db.Exec(
		q, id, userName, firstName, lastName,
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
	var id string
	var err error

	data := dto.CreateUserDTO{
		UserName:  "testUser",
		FirstName: "testFirstName",
		LastName:  "testLastName",
		Email:     "test@email.me",
		Password:  "testPassword",
	}

	ur, db := repo.NewUserSqliteRepo(repo.NewTestDB(), testutil.Fc0)
	id, err = ur.Create(data)

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
    WHERE id = ?`

	var userName string
	var firstName string
	var lastName string
	var email string
	var hashedPassword string
	var createdAt int64
	var updatedAt int64

	err = db.QueryRow(q, id).Scan(
		&userName,
		&firstName,
		&lastName,
		&email,
		&hashedPassword,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		t.Fatalf("%v", err)
	}
	if data.UserName != userName {
		t.Errorf("wanted: %s, got %s", data.UserName, userName)
	}
	if data.FirstName != firstName {
		t.Errorf("wanted: %s, got %s", data.FirstName, firstName)
	}
	if data.LastName != lastName {
		t.Errorf("wanted: %s, got %s", data.LastName, lastName)
	}
	if data.Email != email {
		t.Errorf("wanted: %s, got %s", data.Email, email)
	}
	ok, err := auth.CheckPassword(data.Password, hashedPassword)
	if !ok {
		t.Errorf("password %s did not match hash %s", data.Password, hashedPassword)
	}
	if err != nil {
		t.Errorf(
			"password: %s did not match hash %s, got error %v",
			data.Password,
			hashedPassword,
			err,
		)
	}
	if createdAt != testutil.Fc0.FixedTime.Unix() {
		t.Errorf("want %d, got %d", testutil.Fc0.FixedTime.Unix(), createdAt)
	}
	if updatedAt != testutil.Fc0.FixedTime.Unix() {
		t.Errorf("want %d, got %d", testutil.Fc0.FixedTime.Unix(), updatedAt)
	}
}

func TestDeleteUser(t *testing.T) {
	r, db := repo.NewUserSqliteRepo(repo.NewTestDB(), testutil.Fc0)

	id := "id"
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
		id,
		userName,
		firstName,
		lastName,
		email,
		hashedPassword,
		isActive,
		createdAt,
		updatedAt,
	)
	if err := r.Delete(id); err != nil {
		t.Fatalf("could not delete user due to error: %v", err)
	}
	// TODO: need to check that the user was removed from the db.
}
