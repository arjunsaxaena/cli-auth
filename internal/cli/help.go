package cli

import "fmt"

func showHelp(
	state *State,
) {

	fmt.Println()

	if !state.LoggedIn {

		fmt.Println("Available commands:")

		for _, cmd := range GuestCommands {
			fmt.Println(" -", cmd)
		}

		fmt.Println()
		return
	}

	fmt.Println("Available commands:")

	for _, cmd := range AuthCommands {
		fmt.Println(" -", cmd)
	}

	fmt.Println()
}
