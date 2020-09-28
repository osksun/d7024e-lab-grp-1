package d7024e

import (
	/*
	"strings"
	"fmt"
	"encoding/hex"
	*/
)

// Put command for the cli
func Put() Cmd{
	return Cmd{
		triggers: []string{"put", "p"},
		description: "Upload content",
		usage: "\"put filename content...\", \"p filename content...\"",
		action: func(cli *Cli, args ...string) {
			/*
			filename := args[0]
			content := []byte(strings.Join(args[1:len(args)], " "))
			hash := cli.node.kademlia.Store()
			fmt.Printf("Hash: \"%d\"\nHex : \"%s\"\n"+
				"Stored content \"%s\" in the node's ValueHashtable\n",
				hash, hex.EncodeToString(hash[:]), string(cli.node.Vht().Get([]byte(filename))))
			*/
		},
	}
}