package d7024e

import (
	"bufio"
	"os"
	"fmt"
	"strings"
)

// Cli struct containing bufio reader and list of commands
type Cli struct {
	reader *bufio.Reader
	cmds []Cmd
	node *Node
}

// NewCli Constructor function for Cli class
func NewCli(node *Node) *Cli {
	cli := &Cli{}
	cli.reader = bufio.NewReader(os.Stdin)
	cli.cmds = []Cmd{Help(), Get(), Put(), Exit()} // Add commands here
	cli.node = node
	return cli
}

// Run starts the cli
func (cli *Cli) Run() {
	for {
		fmt.Print(">")
		input, _ := cli.reader.ReadString('\n')
		input = strings.Trim(input,"\r\n")
		if input == "" {
			continue
		}
		args := strings.Split(input, " ")
		unknown := true
		for i := 0; i < len(cli.cmds); i++ {
			if cli.cmds[i].matches(args[0]) {
				cli.cmds[i].action(cli, args[1:]...)
				unknown = false
				break
			}
		}
		if unknown {
			fmt.Printf("Unknown command \"%s\", use command \"help\" for a list of available commands\n", args[0])
		}
	}
}