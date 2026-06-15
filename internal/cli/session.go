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

		fmt.Println(
			"Session expired. Please login again.",
		)

		state.LoggedIn = false
		state.UserID = 0
		state.Username = ""

		return false
	}

	return true
}
