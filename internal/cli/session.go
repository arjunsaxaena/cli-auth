package cli

import (
	"fmt"
	"time"
)

func IsSessionValid(
	state *State,
) bool {

	if !state.LoggedIn {
		return false
	}

	if time.Now().After(
		state.SessionExpiresAt,
	) {

		fmt.Println()
		fmt.Println(
			"Session expired. You have been logged out.",
		)
		fmt.Println()

		state.LoggedIn = false
		state.UserID = 0
		state.Username = ""

		return false
	}

	return true
}
