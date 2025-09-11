package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/users"
	"ssb/internal/models"
	"ssb/internal/schemas"
	"ssb/internal/testutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setup(
	t *testing.T,
	httpMethod string,
	url string,
	body io.Reader,
	us []models.User,
	user models.User,
) (*httptest.ResponseRecorder,
	*testutil.FakeUserRepository) {
	t.Helper()
	req := httptest.NewRequest(httpMethod, url, body)
	w := httptest.NewRecorder()
	ur := testutil.NewFakeUserRepository(us)
	auth := func(request *http.Request) (string, error) { return user.UserName, nil }
	r := users.NewRouter(ur, auth)
	r.ServeHTTP(w, req)
	return w, ur
}

func TestGetUserByUserName(t *testing.T) {
	want := models.User{
		UserName:       "tyler.durdan",
		FirstName:      "tyler",
		LastName:       "durdan",
		Email:          "tyler@paperstreetsoap.com",
		HashedPassword: "secret",
		CreatedAt:      testutil.Fc0.FixedTime.Unix(),
		UpdatedAt:      testutil.Fc0.FixedTime.Unix(),
	}
	w, _ := setup(t, http.MethodGet, "/tyler.durdan", nil, []models.User{want}, want)
	if w.Code != http.StatusOK {
		t.Fatalf("wanted 200, got %d", w.Code)
	}

	var got models.User
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("could not Unmarshal body %s into user model", w.Body.String())
	}

	if !cmp.Equal(want, got) {
		t.Errorf("%v", cmp.Diff(want, got))
	}
}

func TestCreateUser(t *testing.T) {
	admin := models.User{
		UserName: "admin",
	}
	want := models.User{
		UserName:       "tyler.durdan",
		FirstName:      "tyler",
		LastName:       "durdan",
		Email:          "tyler@paperstreetsoap.com",
		HashedPassword: "secret",
		CreatedAt:      testutil.Fc0.FixedTime.Unix(),
		UpdatedAt:      testutil.Fc0.FixedTime.Unix(),
	}

	newUser := schemas.CreateUserDTO{
		UserName:  want.UserName,
		FirstName: want.FirstName,
		LastName:  want.LastName,
		Email:     want.Email,
		Password:  want.HashedPassword,
	}

	data, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("could not marshal dto: %q", newUser)
	}

	w, _ := setup(
		t, http.MethodPost, "/", bytes.NewBuffer(data),
		[]models.User{admin}, admin,
	)
	if w.Code != http.StatusCreated {
		t.Fatalf("wanted %d, got %d", http.StatusCreated, w.Code)
	}

	var got string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("could not Unmarshal body %s into user model", w.Body.String())
	}

	if got != "tyler.durdan" {
		t.Fatalf("wanted tyler.durdan, got %s", got)
	}

}

func TestDeleteUser(t *testing.T) {
	user := models.User{
		UserName:       "tyler.durdan",
		FirstName:      "tyler",
		LastName:       "durdan",
		Email:          "tyler@paperstreetsoap.com",
		HashedPassword: "secret",
	}

	admin := models.User{
		UserName: "admin",
	}

	w := httptest.NewRecorder()
	ur := testutil.NewFakeUserRepository([]models.User{user, admin})
	auth := func(request *http.Request) (string, error) { return "admin", nil }
	r := users.NewRouter(ur, auth)

	url := fmt.Sprintf("/%s", user.UserName)
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("wante %d, got %d", http.StatusNoContent, w.Code)
	}

	if _, ok := ur.UserStore[user.UserName]; ok {
		t.Errorf("User %s is still in the db", user.UserName)
	}
}

func TestUpdateUser(t *testing.T) {
	user := models.User{
		UserName:       "tyler.durdan",
		FirstName:      "tyler",
		LastName:       "durdan",
		Email:          "tyler@paperstreetsoap.com",
		HashedPassword: "secret",
		IsActive:       true,
		CreatedAt:      testutil.Fc0.FixedTime.Unix(),
		UpdatedAt:      testutil.Fc0.FixedTime.Unix(),
	}

	admin := models.User{
		UserName: "admin",
	}

	active := false
	dto := schemas.UpdateUserDTO{
		IsActive: &active,
	}

	data, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("could not marshal dto: %v", dto)
	}

	url := fmt.Sprintf("/%s", user.UserName)
	w, ur := setup(
		t, http.MethodPut, url, bytes.NewBuffer(data),
		[]models.User{user, admin}, admin,
	)

	if w.Code != http.StatusOK {
		t.Fatalf("wanted %d, got %d", http.StatusOK, w.Code)
	}

	updatedUser, ok := ur.UserStore[user.UserName]
	if !ok {
		t.Fatalf("could not find user: %v in user store", user.UserName)
	}

	if active != updatedUser.IsActive {
		t.Errorf(
			"wanted user IsActive: %v, got user IsActive: %v",
			active,
			updatedUser.IsActive,
		)
	}
}
