package cli

import (
	"os"
)

// Exit command for the cli
func Exit() Cmd{
	return Cmd{
		triggers: []string{"exit", "quit", "terminate"},
		description: "Exits the command line interface",
		usage: "\"exit\", \"quit\", \"terminate\"",
		action: func(cli *Cli, args ...string) {
			os.Exit(0)
		},
	}
}
