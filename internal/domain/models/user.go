package models

type User struct {
	UserName       string
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string // bcrypt or argon2
	IsActive       bool
	CreatedAt      int64
	UpdatedAt      int64
}
