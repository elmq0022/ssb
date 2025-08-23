package users_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/users"
	"ssb/internal/models"
	"ssb/internal/testutil"
	"testing"
)

func setup(
	t *testing.T,
	httpMethod string,
	url string,
	body io.Reader,
	us []models.User,
) (*httptest.ResponseRecorder,
	*testutil.FakeUserRepository) {
	t.Helper()
	req := httptest.NewRequest(httpMethod, url, body)
	w := httptest.NewRecorder()
	ur := testutil.NewFakeUserRepository(us)
	r := users.NewRouter(ur)
	r.ServeHTTP(w, req)
	return w, ur
}

func TestGetUserByUserName(t *testing.T) {
	user := models.User{
		UserName:  "tyler.durdan",
		FirstName: "tyler",
		LastName:  "durdan",
		Email:     "tyler@paperstreetsoap.com",
		CreatedAt: testutil.Fc0.FixedTime.Unix(),
		UpdatedAt: testutil.Fc0.FixedTime.Unix(),
	}
	w, _ := setup(t, http.MethodGet, "/tyler.durdan", nil, []models.User{user})
	if w.Code != http.StatusOK {
		t.Fatalf("wanted 200, got %d", w.Code)
	}
}
