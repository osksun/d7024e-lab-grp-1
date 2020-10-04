package d7024e

import (
	"strings"
	"encoding/hex"
)

// Put command for the cli
func Put() Cmd{
	return Cmd{
		triggers: []string{"put", "p"},
		description: "Upload content",
		usage: "\"put filename content...\", \"p filename content...\"",
		action: func(cli *Cli, args ...string) string {
			filename := []byte(args[0])
			content := []byte(strings.Join(args[1:len(args)], " "))
			hash := cli.node.kademlia.Store(filename, content, cli.node.net.storeChannel, cli.node.net.findNodeChannel)
			return "Returned hash (hex): \"" + hex.EncodeToString(hash[:]) + "\"\n"
		},
	}
}