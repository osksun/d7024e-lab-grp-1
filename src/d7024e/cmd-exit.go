package d7024e

import (
	"os"
)

// Exit command for the cli
func Exit() Cmd{
	return Cmd{
		triggers: []string{"exit", "quit", "terminate", "e", "q", "t"},
		description: "Terminates the node",
		usage: "\"exit\", \"quit\", \"terminate\"",
		action: func(cli *Cli, args ...string) string {
			os.Exit(0)
			return "" // This will never be executed since we're exiting before
		},
	}
}
