package commands

import (
	"fmt"
	"os"
)

func HandleUser(args []string) {
	switch args[0] {
	case "create":
		fmt.Println("hello from create user")
	case "delete":
		fmt.Println("hello from delete user")
	default:
		fmt.Println("expected one of 'create' or 'delete'")
		os.Exit(1)
	}
}

func HandleCreateUser(client HTTPClient) {
}
