package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadSetupMethod() string {

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Println()
		fmt.Println("Choose setup method:")
		fmt.Println("1. QR Code")
		fmt.Println("2. Secret Key")

		fmt.Print("Selection: ")

		choice, _ := reader.ReadString('\n')

		choice = strings.TrimSpace(
			choice,
		)

		if choice == "1" ||
			choice == "2" {

			return choice
		}

		fmt.Println(
			"Invalid selection.",
		)
	}
}
