package cli

import (
	"bufio"
	"cli-auth/internal/auth"
	"fmt"
	"os"
	"strings"
	"time"
)

func handleLogin(
	repo *auth.Repository,
	state *State,
) {

	if state.LoggedIn {
		fmt.Println("Already logged in.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")

	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	password, err := ReadPassword("Password: ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	user, err := auth.Login(
		repo,
		username,
		password,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	state.LoggedIn = true
	state.UserID = user.ID
	state.Username = user.Username
	state.SessionExpiresAt = time.Now().Add(auth.SessionDuration)

	fmt.Println()
	fmt.Println("Login successful.")
	fmt.Println()

	showUserDetails(user, state)
}
