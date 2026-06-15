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

	user, err := repo.GetUserByID(
		state.UserID,
	)

	if err != nil {
		fmt.Println("Unable to load user.")
		return
	}

	if !user.MFAEnabled {
		fmt.Println(
			"2FA is already disabled.",
		)
		return
	}

	if user.TOTPSecret == nil {
		fmt.Println(
			"Missing TOTP secret.",
		)
		return
	}

	code := ReadTOTP()

	if !auth.ValidateTOTP(
		code,
		*user.TOTPSecret,
	) {

		fmt.Println(
			"Invalid TOTP code.",
		)

		return
	}

	err = repo.DisableMFA(
		user.ID,
	)

	if err != nil {
		fmt.Println(
			"Error:",
			err,
		)

		return
	}

	fmt.Println(
		"2FA disabled successfully.",
	)
}
