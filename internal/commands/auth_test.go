package commands_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	cmd "ssb/internal/commands"
	tu "ssb/internal/commands/testUtils"
	"ssb/internal/commands/utils"
	"ssb/internal/schemas"
	"testing"
)

func NewPasswordFunc(password string) cmd.PasswordFunc {
	return func() (string, error) {
		return password, nil
	}
}

func TestHandleLogin_Success(t *testing.T) {
	empty := schemas.JsonToken{}
	tu.SetJWTToken(t, empty)

	want := schemas.JsonToken{
		Token: "test-token",
	}
	body, _ := json.Marshal(want)
	client := &tu.FakeClient{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(body)),
		},
		Err: nil,
	}

	pf := NewPasswordFunc("password")

	if err := cmd.HandleLogin(pf, client); err != nil {
		t.Fatalf("expected no error, bot %v", err)
	}

	got := utils.MustReadJWTToken()

	if want.Token != got.Token {
		t.Fatalf("want %s, got %s", want.Token, got.Token)
	}
}
