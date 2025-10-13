package commands

import "fmt"

func HandleConfig(args []string) {
	switch args[0] {
	case "set-server":
		fmt.Println("hello from set-server")
	case "set-user":
		fmt.Println("hello from set-user")
	default:
		fmt.Println("expected one of 'set-server' or 'set-user'")
	}
}
