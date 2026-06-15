package cli

import (
	"bufio"
	"cli-auth/internal/auth"
	"cli-auth/internal/config"
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

	user, err := auth.AuthenticatePassword(
		repo,
		username,
		password,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if user.MFAEnabled {

		if user.TOTPSecret == nil {
			fmt.Println("User has MFA enabled but no secret configured.")
			return
		}

		code := ReadTOTP()

		if !auth.ValidateTOTP(
			code,
			*user.TOTPSecret,
		) {
			fmt.Println("Invalid TOTP code.")
			return
		}
	}

	err = auth.CompleteLogin(
		repo,
		user.ID,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	user, err = repo.GetUserByID(
		user.ID,
	)

	if err != nil {
		fmt.Println(
			"Error loading updated user:",
			err,
		)
		return
	}

	state.LoggedIn = true
	state.UserID = user.ID
	state.Username = user.Username
	state.SessionExpiresAt = time.Now().
		Add(config.App.SessionDuration)

	fmt.Println()
	fmt.Println("Login successful.")
	fmt.Println()

	showUserDetails(
		user,
		state,
	)
}
