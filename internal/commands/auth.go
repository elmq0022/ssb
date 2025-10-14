package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"ssb/internal/commands/utils"
	"ssb/internal/schemas"
	"strings"

	"golang.org/x/term"
)

type PasswordFunc func() (string, error)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func HandleAuth(args []string) {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, "expected one of: 'login'")
		os.Exit(2)
	}

	switch args[0] {
	case "login":
		if err := HandleLogin(GetPasswordFromStdin, http.DefaultClient); err != nil {
			fmt.Fprint(os.Stderr, "login failed:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("expected one of: 'login'")
		os.Exit(2)
	}
}

func GetPasswordFromStdin() (string, error) {
	fmt.Fprint(os.Stdout, "Password: ")
	bytePw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePw)), nil
}

func HandleLogin(pf PasswordFunc, client HTTPClient) error {
	password, err := pf()
	if err != nil {
		return fmt.Errorf("could not read password: %w", err)
	}

	cfg := utils.MustReadConfig()
	server := cfg.URL
	user := cfg.Username

	data := schemas.LoginRequest{
		Username: user,
		Password: password,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		server+"/auth/login",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send login request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", resp.Status)
	}
	defer resp.Body.Close()

	var jwtToken schemas.JsonToken
	if err := json.NewDecoder(resp.Body).Decode(&jwtToken); err != nil {
		return fmt.Errorf("could not decode token: %w", err)
	}

	utils.MustSetJWTToken(jwtToken)

	return nil
}
