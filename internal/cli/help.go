package cli

import "fmt"

func showHelp(
	state *State,
) {

	fmt.Println()

	if !state.LoggedIn {

		fmt.Println("Available Commands")
		fmt.Println("------------------")
		fmt.Println("register      Register a new user")
		fmt.Println("login         Login to your account")
		fmt.Println("help          Show available commands")
		fmt.Println("exit          Exit application")

		fmt.Println()
		return
	}

	fmt.Println("Available Commands")
	fmt.Println("------------------")
	fmt.Println("whoami        Show current user details")
	fmt.Println("enable-2fa    Enable TOTP authentication")
	fmt.Println("disable-2fa   Disable TOTP authentication")
	fmt.Println("logout        Logout current user")
	fmt.Println("help          Show available commands")
	fmt.Println("exit          Exit application")

	fmt.Println()
}
