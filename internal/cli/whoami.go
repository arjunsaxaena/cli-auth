package cli

import (
	"cli-auth/internal/auth"
	"fmt"
)

func handleWhoAmI(
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
		fmt.Println(
			"Unable to fetch user details.",
		)
		return
	}

	showUserDetails(
		user,
		state,
	)
}
