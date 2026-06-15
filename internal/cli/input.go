package cli

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	password, err := term.ReadPassword(
		int(syscall.Stdin),
	)

	fmt.Println()

	if err != nil {
		return "", err
	}

	return string(password), nil
}
