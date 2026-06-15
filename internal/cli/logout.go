package cli

import (
	"fmt"
	"time"
)

func handleLogout(
	state *State,
) {

	if !state.LoggedIn {
		fmt.Println("Not logged in.")
		return
	}

	state.LoggedIn = false
	state.UserID = 0
	state.Username = ""
	state.SessionExpiresAt = time.Time{}

	fmt.Println("Logged out successfully.")
}
