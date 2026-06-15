package cli

import (
	"cli-auth/internal/auth"
	"fmt"
)

func handleEnable2FA(
	repo *auth.Repository,
	state *State,
) {

	if !IsSessionValid(state) {
		return
	}

	user, err := repo.GetUserByID(
		state.UserID,
	)

	if err != nil {
		fmt.Println("Unable to load user.")
		return
	}

	if user.MFAEnabled {
		fmt.Println("2FA already enabled.")
		return
	}

	key, err := auth.GenerateTOTPKey(
		user.Username,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	choice := ReadSetupMethod()

	err = repo.EnableMFA(
		user.ID,
		key.Secret(),
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch choice {

	case "1":

		path, err := auth.GenerateQRCode(
			user.Username,
			key.URL(),
		)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println()
		fmt.Println(
			"2FA enabled successfully.",
		)

		fmt.Printf(
			"\nQR Code saved to: %s\n",
			path,
		)

		fmt.Println(
			"Scan it with Google Authenticator.",
		)

	case "2":

		fmt.Println()
		fmt.Println(
			"2FA enabled successfully.",
		)

		fmt.Println(
			"\nAdd this secret to your authenticator app:",
		)

		fmt.Println()
		fmt.Println(
			key.Secret(),
		)
	}
}
