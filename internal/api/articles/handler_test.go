package articles_test

import (
	"net/http"
	"net/http/httptest"
	"ssb/internal/api/articles"
	"ssb/internal/domain/models"
	"testing"
)

func TestGetArticles(t *testing.T) {
	articles := []models.Article{}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r := articles.NewRouter()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	expectedBody := "\"\"\n"

	if w.Body.String() != expectedBody {
		t.Errorf("unexpected body: got %q, want %q", w.Body.String(), expectedBody)
	}
}
