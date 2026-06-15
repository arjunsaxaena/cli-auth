package cli

import (
	"fmt"
	"strings"

	"cli-auth/internal/auth"

	"github.com/chzyer/readline"
)

func Start(
	repo *auth.Repository,
) error {

	state := &State{}

	completer := readline.NewPrefixCompleter(
		readline.PcItem("register"),
		readline.PcItem("login"),
		readline.PcItem("help"),
		readline.PcItem("exit"),
		readline.PcItem("whoami"),
		readline.PcItem("enable-2fa"),
		readline.PcItem("disable-2fa"),
		readline.PcItem("logout"),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/fortress-history",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		return err
	}

	defer rl.Close()

	for {

		line, err := rl.Readline()
		if err != nil {
			break
		}

		cmd := strings.TrimSpace(line)

		switch cmd {

		case "help":
			showHelp(state)

		case "register":
			handleRegister(repo)

		case "exit":
			fmt.Println("Goodbye.")
			return nil

		default:
			fmt.Println("Unknown command. Type 'help'.")
		}
	}

	return nil
}
