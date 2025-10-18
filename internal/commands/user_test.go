package commands_test

import (
	"bytes"
	"io"
	"net/http"
	cmd "ssb/internal/commands"
	tu "ssb/internal/commands/testUtils"
	"ssb/internal/schemas"
	"testing"
)

func TestHandleCreateUser(t *testing.T) {
	tu.SetJWTToken(t, schemas.JsonToken{
		Token: "test-token",
	})
	fc := &tu.FakeClient{
		Resp: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader(nil)),
		},
		Err: nil,
	}
	data := schemas.CreateUserDTO{}
	if err := cmd.HandleCreateUser(data, fc); err != nil {
		t.Fatalf("create user returned error: %q", err)
	}
}
