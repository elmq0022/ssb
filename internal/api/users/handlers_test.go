package users_test

import (
	"io"
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

func TestGetUser(t *testing.T) {
	w, ur := setup()
}
