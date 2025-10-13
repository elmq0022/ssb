package commands_test

import (
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

func TestHandleLogin(t *testing.T) {
	cfg := utils.CLIConfig{
		URL:      "",
		Username: "",
	}

	jwt := schemas.JsonToken{
		Token: "bad-token",
	}

	tu.SetConfig(t, cfg)
	tu.SetJWTToken(t, jwt)

	pf := NewPasswordFunc("password")
	cmd.HandleLogin(pf)

}
