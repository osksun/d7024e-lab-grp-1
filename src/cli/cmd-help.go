package cli

import (
	"fmt"
	"strings"
)

// Help command for the cli
func Help() Cmd{
	return Cmd{
		triggers: []string{"help", "h"},
		description: "Shows list of available commands",
		usage: "\"help\", \"h\"",
		action: func(cli *Cli, args ...string) {
			fmt.Printf("%-30s| %-40s| %-40s\n", "Trigger(s)", "Description", "Usage")
			for i := 0; i < len(cli.cmds); i++ {
				fmt.Printf("%-30s| %-40s| %-40s\n",
					strings.Join(cli.cmds[i].triggers, ", "),
					cli.cmds[i].description,
					cli.cmds[i].usage,
				)
			}
		},
	}
}