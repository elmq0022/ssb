package testutil

import (
	"errors"
	"ssb/internal/models"
	"ssb/internal/repo"
	"ssb/internal/schemas"
)

// supposedly helps by providing a compile time check of the interface
var _ repo.UserRepository = (*FakeUserRepository)(nil)

type FakeUserRepository struct {
	UserStore map[string]models.User
}

func NewFakeUserRepository(users []models.User) *FakeUserRepository {
	us := make(map[string]models.User)

	for _, user := range users {
		us[user.UserName] = user
	}

	f := FakeUserRepository{
		UserStore: us,
	}

	return &f
}

func (f *FakeUserRepository) GetByUserName(username string) (models.User, error) {
	user, ok := f.UserStore[username]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (f *FakeUserRepository) Create(data schemas.CreateUserDTO) (string, error) {
	return "", errors.New("NotImplemented")
}

func (f *FakeUserRepository) Update(userName string, data schemas.UpdateUserDTO) error {
	return errors.New("NotImplemented")
}

func (f *FakeUserRepository) Delete(userName string) error {
	return errors.New("NotImplemented")
}
