package commands

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func HandleAuth(args []string) {
	switch args[0] {
	case "login":
		HandleLogin(GetPasswordFromStdin)
	default:
		fmt.Println("expected one of 'login'")
		os.Exit(1)
	}
}

type PasswordFunc func() (string, error)

func GetPasswordFromStdin() (string, error) {
	fmt.Fprint(os.Stdout, "Password: ")
	bytePw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	return string(bytePw), nil
}

func HandleLogin(pf PasswordFunc) {
}
