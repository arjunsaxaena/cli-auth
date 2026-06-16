package cli

import (
	"fmt"
	"strings"

	"cli-auth/internal/auth"

	"github.com/chzyer/readline"
)

type StateCompleter struct {
	state *State
	guest *readline.PrefixCompleter
	auth  *readline.PrefixCompleter
}

func (sc *StateCompleter) Do(line []rune, pos int) ([][]rune, int) {
	if sc.state.LoggedIn {
		return sc.auth.Do(line, pos)
	}
	return sc.guest.Do(line, pos)
}

func Start(
	repo *auth.Repository,
) error {

	state := &State{}

	guestCompleter := readline.NewPrefixCompleter(
		readline.PcItem("register"),
		readline.PcItem("login"),
		readline.PcItem("help"),
		readline.PcItem("exit"),
	)

	authCompleter := readline.NewPrefixCompleter(
		readline.PcItem("whoami"),
		readline.PcItem("enable-2fa"),
		readline.PcItem("disable-2fa"),
		readline.PcItem("logout"),
		readline.PcItem("help"),
		readline.PcItem("exit"),
	)

	completer := &StateCompleter{
		state: state,
		guest: guestCompleter,
		auth:  authCompleter,
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "guest > ",
		HistoryFile:     "/tmp/cli-auth-history",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		return err
	}

	defer rl.Close()

	for {
		if state.LoggedIn {
			rl.SetPrompt(fmt.Sprintf("\033[1;32m[%s]\033[0m > ", state.Username))
		} else {
			rl.SetPrompt("\033[1;34mguest\033[0m > ")
		}

		line, err := rl.Readline()
		if err != nil {
			break
		}

		cmd := strings.TrimSpace(line)
		if cmd == "" {
			continue
		}

		if state.LoggedIn &&
			!IsSessionValid(state) {
			continue
		}

		switch cmd {

		case "help":
			showHelp(state)

		case "register":
			if state.LoggedIn {
				fmt.Println("Error: You are already logged in.")
			} else {
				handleRegister(repo)
			}

		case "exit":
			fmt.Println("Goodbye.")
			return nil

		case "login":
			if state.LoggedIn {
				fmt.Println("Error: You are already logged in.")
			} else {
				handleLogin(
					repo,
					state,
				)
			}

		case "whoami":
			if !state.LoggedIn {
				fmt.Println("Error: You must be logged in to use this command.")
			} else {
				handleWhoAmI(
					repo,
					state,
				)
			}

		case "enable-2fa":
			if !state.LoggedIn {
				fmt.Println("Error: You must be logged in to use this command.")
			} else {
				handleEnable2FA(
					repo,
					state,
				)
			}

		case "disable-2fa":
			if !state.LoggedIn {
				fmt.Println("Error: You must be logged in to use this command.")
			} else {
				handleDisable2FA(
					repo,
					state,
				)
			}

		case "logout":
			if !state.LoggedIn {
				fmt.Println("Error: You must be logged in to use this command.")
			} else {
				handleLogout(state)
			}

		default:
			fmt.Println("Unknown command. Type 'help'.")
		}
	}

	return nil
}
