package d7024e

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
		usage: "\"put filename content...\", \"p filename content...\"",
		action: func(cli *Cli, args ...string) {
			filename := []byte(args[0])
			content := []byte(strings.Join(args[1:len(args)], " "))
			hash := cli.node.kademlia.Store(filename, content)
			fmt.Printf("Returned hash (hex): \"%s\"\n", hex.EncodeToString(hash[:]))
		},
	}
}