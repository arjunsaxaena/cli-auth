package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cli-auth/internal/auth"
)

func handleRegister(
	repo *auth.Repository,
) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")

	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	password, err := ReadPassword("Password: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = auth.Register(
		repo,
		username,
		password,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("User created successfully.")
}
