package d7024e

import (
	"fmt"
	"encoding/hex"
)

// Get command for the cli
func Get() Cmd{
	return Cmd{
		triggers: []string{"get", "g"},
		description: "Get content of a file",
		usage: "\"get hashOfFileName\", \"g hashOfFileName\"",
		action: func(cli *Cli, args ...string) {
			filenameHashSlice, _ := hex.DecodeString(args[0])
			var filenameHash [HashSize]byte
			copy(filenameHash[:], filenameHashSlice)
			data := cli.node.kademlia.LookupData(filenameHash)
			fmt.Printf("Returned content: \"%s\"\n", string(data))
		},
	}
}