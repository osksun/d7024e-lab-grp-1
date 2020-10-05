package d7024e

import (
	"strings"
)

// Cmd struct containing name, aliases, description, usage and action
type Cmd struct {
	triggers []string
	description string
	usage string
	action func(cli *Cli, args ...string) string
}

// matches checks if given input matches any trigger of the command
func (cmd *Cmd) matches(input string) bool {
	input = strings.ToLower(input)
	for i := 0; i < len(cmd.triggers); i++ {
		if strings.Compare(input, cmd.triggers[i]) == 0 {
			return true
		}
	}
	return false
}