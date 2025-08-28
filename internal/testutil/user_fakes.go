package testutil

import (
	"errors"
	"ssb/internal/models"
	"ssb/internal/pkg/auth"
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
	if _, ok := f.UserStore[data.UserName]; ok {
		return "", errors.New("user already exists")
	}
	newUser := models.User{
		UserName:       data.UserName,
		FirstName:      data.FirstName,
		LastName:       data.LastName,
		Email:          data.Email,
		HashedPassword: auth.HashPassword(data.Password),
		CreatedAt:      Fc0.FixedTime.Unix(),
		UpdatedAt:      Fc0.FixedTime.Unix(),
	}
	f.UserStore[newUser.UserName] = newUser
	return newUser.UserName, nil
}

func (f *FakeUserRepository) Update(userName string, data schemas.UpdateUserDTO) error {
	user, ok := f.UserStore[userName]
	if !ok {
		return errors.New("User does not exist")
	}

	if data.IsActive != nil {
		user.IsActive = *data.IsActive
	}

	f.UserStore[userName] = user
	return nil
}

func (f *FakeUserRepository) Delete(userName string) error {
	delete(f.UserStore, userName)
	return nil
}
