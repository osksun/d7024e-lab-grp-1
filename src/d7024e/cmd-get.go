package d7024e

import (
	"fmt"
)

// Get command for the cli
func Get() Cmd{
	return Cmd{
		triggers: []string{"get", "g"},
		description: "Get content of a file",
		usage: "\"get filename\", \"g filename\"",
		action: func(cli *Cli, args ...string) {
			filename := args[0]
			fmt.Printf("Returned content: \"%s\"\n", string(cli.node.Vht().Get([]byte(filename))))
		},
	}
}