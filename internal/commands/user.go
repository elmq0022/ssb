package commands

import (
	"fmt"
	"net/http"
	"os"
	"ssb/internal/commands/utils"
	"ssb/internal/schemas"
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

func HandleCreateUser(userData schemas.CreateUserDTO, client HTTPClient) error {
	ep, err := utils.BuildEndpoint("users")
	if err != nil {
		return fmt.Errorf("could not build url for user endpoint: %w", err)
	}

	req, err := utils.NewRequestBuilder(http.MethodPost, ep).
		WithJSON(userData).
		WithAuth().
		Build()
	if err != nil {
		return fmt.Errorf("could not build request for create user: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not execute request: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("did not receive status code %d: %w", resp.StatusCode, err)
	}

	return nil
}
