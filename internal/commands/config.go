package commands

import (
	"bufio"
	"fmt"
	"os"
	"ssb/internal/commands/utils"
)

func HandleConfig(args []string) {
	switch args[0] {
	case "init":
		handleInit()
	case "set-server":
		fmt.Println("hello from set-server")
	case "set-user":
		fmt.Println("hello from set-user")
	default:
		fmt.Println("expected one of 'set-server' or 'set-user'")
	}
}

func handleInit(){
	cfg := utils.CLIConfig{}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter bfs url: ")
	scanner.Scan()
	cfg.URL = scanner.Text()

	fmt.Print("Enter your username: ")
	scanner.Scan()
	cfg.Username = scanner.Text()

	utils.MustSetConfig(&cfg)
}
