package cli

import (
	"cli-auth/internal/auth"
	"fmt"
)

func handleDisable2FA(
	repo *auth.Repository,
	state *State,
) {

	if !IsSessionValid(state) {
		return
	}

	err := repo.DisableMFA(
		state.UserID,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("2FA disabled.")
}
