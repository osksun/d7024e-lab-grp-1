package d7024e

import (
	"bufio"
	"fmt"
	"strings"
	"io"
)

// Cli struct containing bufio reader and list of commands
type Cli struct {
	reader *bufio.Reader
	cmds []Cmd
	node *Node
}

// NewCli Constructor function for Cli class
func NewCli(node *Node, rd io.Reader) *Cli {
	cli := &Cli{}
	cli.reader = bufio.NewReader(rd)
	cli.cmds = []Cmd{Help(), Get(), Put(), Exit(), Contacts()} // Add commands here
	cli.node = node
	return cli
}

// Run starts the cli
func (cli *Cli) Run(runOnce bool, verbose bool) int {
	if runOnce {
		cli.handleInput(cli.getInput(), verbose)
		return 0
	}
	for {
		cli.handleInput(cli.getInput(), verbose)
	}
}

func (cli *Cli) getInput() string {
	fmt.Print(">")
	input, _ := cli.reader.ReadString('\n')
	input = strings.Trim(input,"\n")
	input = strings.Trim(input,"\r")
	return input
}

func (cli *Cli) handleInput(input string, verbose bool) {
	if input != "" {
		args := strings.Split(input, " ")
		unknown := true
		for i := 0; i < len(cli.cmds); i++ {
			if cli.cmds[i].matches(args[0]) {
				result := cli.cmds[i].action(cli, args[1:]...)
				if verbose {
					fmt.Println(result)
				}
				unknown = false
				break
			}
		}
		if unknown && verbose {
			fmt.Printf("Unknown command \"%s\", use command \"help\" for a list of available commands\n", args[0])
		}
	}
}