package main

import (
	"fmt"
	"os"
	cmd "ssb/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected one of 'auth', 'user', or 'article'")
		os.Exit(1)
	}

	fmt.Println(os.Args)

	switch os.Args[1] {
	case "auth":
		cmd.HandleAuth(os.Args[2:])
	case "user":
		// TODO: handle user
	case "article":
		// TODO: handle article
	default:
		fmt.Println("expected one of 'auth', 'user', or 'article'")
		os.Exit(1)
	}
}
