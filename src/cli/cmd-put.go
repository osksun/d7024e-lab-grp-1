package cli

import (
	"strings"
	"fmt"
	"encoding/hex"
)

// Put command for the cli
func Put() Cmd{
	return Cmd{
		triggers: []string{"put", "p"},
		description: "Upload content",
		usage: "\"put CONTENT HERE\", \"p CONTENT HERE\"",
		action: func(cli *Cli, args ...string) {
			filename := args[0]
			content := []byte(strings.Join(args[1:len(args)], " "))
			hash := cli.node.Vht().Put([]byte(filename), []byte(content))
			fmt.Printf("Hash: \"%d\"\n(Hex): \"%s\"\n"+
				"Stored data: \"%s\" in the node's ValueHashtable\n",
				hash, hex.EncodeToString(hash[:]), string(cli.node.Vht().Get([]byte(filename))))
		},
	}
}