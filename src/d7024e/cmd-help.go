package d7024e

import (
	"strings"
)

// Help command for the cli
func Help() Cmd{
	return Cmd{
		triggers: []string{"help", "h"},
		description: "Shows list of available commands",
		usage: "\"help\", \"h\"",
		action: func(cli *Cli, args ...string) string {
			result := "Trigger(s), Description, Usage"
			for i := 0; i < len(cli.cmds); i++ {
				result += "\n" + strings.Join(cli.cmds[i].triggers, ", ") + " | " + cli.cmds[i].description + " | " + cli.cmds[i].usage
			}
			return result
		},
	}
}