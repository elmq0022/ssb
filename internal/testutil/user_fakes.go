package testutil

import (
	"errors"
	"ssb/internal/models"
	"ssb/internal/schemas"
)

type FakeUserRepository struct {
	UserStore map[string]models.User
}

func NewFakeUserRepository(users []models.User) FakeUserRepository {
	us := make(map[string]models.User)

	for _, user := range users {
		us[user.UserName] = user
	}

	return FakeUserRepository{
		UserStore: us,
	}
}

func (f *FakeUserRepository) GetByUserName(username string) (models.User, error) {
	return models.User{}, errors.New("NotImplemented")
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
