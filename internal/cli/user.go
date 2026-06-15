package cli

import (
	"fmt"

	"cli-auth/internal/models"
)

func showUserDetails(
	user *models.User,
	state *State,
) {

	fmt.Println("User Details")
	fmt.Println("------------")

	fmt.Printf(
		"Username          : %s\n",
		user.Username,
	)

	fmt.Printf(
		"Registered        : %s\n",
		user.CreatedAt.Format(
			"2006-01-02 15:04:05",
		),
	)

	if user.MFAEnabled {
		fmt.Println(
			"MFA Status        : Enabled",
		)
	} else {
		fmt.Println(
			"MFA Status        : Disabled",
		)
	}

	fmt.Printf(
		"Session Expires   : %s\n",
		state.SessionExpiresAt.Format(
			"2006-01-02 15:04:05",
		),
	)

	if user.LastLogin != nil {

		fmt.Printf(
			"Last Login        : %s\n",
			user.LastLogin.Format(
				"2006-01-02 15:04:05",
			),
		)

	} else {

		fmt.Println(
			"Last Login        : Never",
		)
	}

	fmt.Println()
}
