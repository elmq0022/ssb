package commands

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"ssb/internal/commands/utils"
	"ssb/internal/schemas"
)

func HandleUser(args []string) {
	switch args[0] {
	case "create":
		dto := schemas.CreateUserDTO{}
		dto.UserName = DefaultPrompt("enter username: ")
		dto.FirstName = DefaultPrompt("enter first name: ")
		dto.LastName = DefaultPrompt("enter last name: ")
		dto.Email = DefaultPrompt("enter email: ")

		password, err := DefaultReadPasswordTwice()
		if err != nil {
			fmt.Printf("could not set password: %q", err)
			os.Exit(1)
		}
		dto.Password = password

		err = HandleCreateUser(dto, &http.Client{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "delete":
		fmt.Println("Not Implemented")
	default:
		fmt.Println("expected one of 'create' or 'delete'")
		os.Exit(1)
	}
}

func Prompt(r io.Reader, w io.Writer, prompt string) string {
	fmt.Fprint(w, prompt)
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	return scanner.Text()
}

func DefaultPrompt(prompt string) string {
	return Prompt(os.Stdin, os.Stderr, prompt)
}

func ReadPasswordTwice(r io.Reader, w io.Writer) (string, error) {
	p1 := Prompt(r, w, "enter your password: ")
	p2 := Prompt(r, w, "enter your password again: ")
	if p1 != p2 {
		return "", fmt.Errorf("passwords do not match")
	}
	return p1, nil
}

func DefaultReadPasswordTwice() (string, error) {
	return ReadPasswordTwice(os.Stdin, os.Stderr)
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
