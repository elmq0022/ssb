package testutils

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"ssb/internal/commands/utils"
	"ssb/internal/schemas"
	"testing"
)

func failOnErr(t *testing.T, e error) {
	t.Helper()
	if e != nil {
		t.Fatalf("failed with error: %q", e)
	}
}

func SetConfig(t *testing.T, cfg utils.CLIConfig) {
	t.Helper()

	f := filepath.Join(t.TempDir(), "config.json")

	data, err := json.MarshalIndent(cfg, "", " ")
	failOnErr(t, err)

	failOnErr(t, os.WriteFile(f, data, 0o600))

	utils.ConfigFilePath = f
}

func SetJWTToken(t *testing.T, token schemas.JsonToken) {
	t.Helper()

	f := filepath.Join(t.TempDir(), "token.json")

	data, err := json.MarshalIndent(token, "", " ")
	failOnErr(t, err)

	failOnErr(t, os.WriteFile(f, data, 0o600))

	utils.JWTFilePath = f
}

type FakeClient struct {
	Resp *http.Response
	Err  error
}

func (f *FakeClient) Do(req *http.Request) (*http.Response, error) {
	return f.Resp, f.Err
}
