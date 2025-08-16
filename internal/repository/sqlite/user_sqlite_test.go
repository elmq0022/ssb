package repo_test

import (
	"database/sql"
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

	// TODO: test hashed password

	if createdAt != testutil.Fc0.FixedTime.Unix() {
		t.Errorf("want %d, got %d", testutil.Fc0.FixedTime.Unix(), createdAt)
	}

	if updatedAt != testutil.Fc0.FixedTime.Unix() {
		t.Errorf("want %d, got %d", testutil.Fc0.FixedTime.Unix(), updatedAt)
	}
}
