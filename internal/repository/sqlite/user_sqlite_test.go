package repo_test

import (
	"database/sql"
	"ssb/internal/dto"
	"ssb/internal/repository"
	"ssb/internal/testutil"
	"ssb/internal/timeutil"
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

func TestCreateUser(t *testing.T){
    var id string
    var err error

    data := dto.CreateUserDTO{
        UserName: "testUser",
        FirstName: "testFirstName",
        LastName: "testLastName",
        Email: "test@email.me",
        Password: "testPassword",
    }
    db := NewUserTestDB(t)
    r := NewUserSqliteRepo(db, testutil.Fc0)
    id, err = r.CreateUser(data)

    if err != nil {
        t.Fatalf("%v", err)
    }

    q :=`
    SELECT 
      userName, 
      firstName, 
      lastName, 
      email, 
      hashedPassword, 
      createdAt,
      updatedAt 
    FROM USERS 
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
    // TODO: test the actual values created
}
