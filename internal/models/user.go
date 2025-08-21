package models

type User struct {
	UserName       string `db:"user_name"`
	FirstName      string `db:"first_name"`
	LastName       string `db:"last_name"`
	Email          string `db:"email"`
	HashedPassword string `db:"hashed_password"`
	IsActive       bool   `db:"is_active"`
	CreatedAt      int64  `db:"created_at"`
	UpdatedAt      int64  `db:"updated_at"`
}
