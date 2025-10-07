package commands

import (
	"fmt"
	"os"
)

func HandleAuth(args []string) {
	switch args[0] {
	case "login":
		fmt.Println("hello from login")
	default:
		fmt.Println("expected one of 'login'")
		os.Exit(1)
	}
}
