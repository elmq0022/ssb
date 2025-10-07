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
		cmd.HandleUser(os.Args[2:])
	case "article":
		cmd.HandleArticle(os.Args[2:])
	default:
		fmt.Println("expected one of 'auth', 'user', or 'article'")
		os.Exit(1)
	}
}
