package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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
