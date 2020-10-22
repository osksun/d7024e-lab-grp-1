package d7024e

// Get command for the cli
func Get() Cmd{
	return Cmd{
		triggers: []string{"get", "g"},
		description: "Get content of a file",
		usage: "\"get fileName\", \"g fileName\"",
		action: func(cli *Cli, args ...string) string {
			filenameHash := Hash([]byte(args[0]))
			data := cli.node.kademlia.LookupData(filenameHash, cli.node.net.findDataChannel)
			return "Returned content: \"" + string(data) + "\""
		},
	}
}