package commands_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	cmd "ssb/internal/commands"
	"ssb/internal/schemas"
	"testing"
)

func NewPasswordFunc(password string) cmd.PasswordFunc {
	return func() (string, error) {
		return password, nil
	}
}

type fakeClient struct {
	resp *http.Response
	err  error
}

func (f *fakeClient) Do(req *http.Request) (*http.Response, error) {
	return f.resp, f.err
}

func TestHandleLogin_Success(t *testing.T) {
	token := schemas.JsonToken{
		Token: "test-token",
	}
	body, _ := json.Marshal(token)
	client := &fakeClient{
		resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		},
	}

	pf := NewPasswordFunc("password")

	if err := cmd.HandleLogin(pf, client); err != nil {
		t.Fatalf("expected no error, bot %v", err)
	}

}
