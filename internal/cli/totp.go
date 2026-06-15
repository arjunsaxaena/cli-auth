package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadTOTP() string {

	reader :=
		bufio.NewReader(os.Stdin)

	fmt.Print(
		"TOTP Code: ",
	)

	code, _ :=
		reader.ReadString('\n')

	return strings.TrimSpace(
		code,
	)
}
