package d7024e

import (
	"testing"
	"os"
	"os/exec"
	"fmt"
	"strings"
)

func TestNewCli(t *testing.T) {
	node := NewNode("localhost:0000")
	cli := NewCli(node, os.Stdin)
	cliType := fmt.Sprintf("%T", cli)
	if cliType != "*d7024e.Cli" {
		t.Errorf("The returned object is not of type \"*d7024e.Cli\"")
	}
	readerType := fmt.Sprintf("%T", cli.reader)
	if readerType != "*bufio.Reader" {
		t.Errorf("The returned object.reader is not of type \"*bufio.Reader\"")
	}
	cmdsType := fmt.Sprintf("%T", cli.cmds)
	if cmdsType != "[]d7024e.Cmd" {
		t.Errorf("The returned object.cmds is not of type \"[]d7024e.Cmd\"")
	}
	nodeType := fmt.Sprintf("%T", cli.node)
	if nodeType != "*d7024e.Node" {
		t.Errorf("The returned object.cmds is not of type \"*d7024e.Node\"")
	}
}

func TestRun(t *testing.T) {
	node := NewNode("localhost:0000")

	cli1 := NewCli(node, strings.NewReader("help\n"))
	statusCode1 := cli1.Run(true, false)
	if statusCode1 != 0 {
		t.Errorf("Expected status code \"0\" but got %d", statusCode1)
	}

	cli2 := NewCli(node, strings.NewReader("invalid command\nexit\n"))
	statusCode2 := cli2.Run(true, false)
	if statusCode2 != 0 {
		t.Errorf("Expected status code \"0\" but got %d", statusCode2)
	}

	cli3 := NewCli(node, strings.NewReader("help\nexit\n"))
	if os.Getenv("BE_EXIT") == "1" {
		cli3.Run(true, false)
		return
    }
    cmd1 := exec.Command(os.Args[0], "-test.run=TestExit")
    cmd1.Env = append(os.Environ(), "BE_EXIT=1")
    err1 := cmd1.Run()
	e1, exitedIncorrectly1 := err1.(*exec.ExitError)
	if exitedIncorrectly1 {
		t.Errorf("Expected exit status message message \"exit status 0\" but got \"%s\"\n", e1.Error())
	}

	cli4 := NewCli(node, strings.NewReader("invalid command\nexit\n"))
	if os.Getenv("BE_EXIT") == "1" {
		cli4.Run(true, false)
		return
    }
    cmd2 := exec.Command(os.Args[0], "-test.run=TestExit")
    cmd2.Env = append(os.Environ(), "BE_EXIT=1")
    err2 := cmd2.Run()
	e2, exitedIncorrectly2 := err2.(*exec.ExitError)
	if exitedIncorrectly2 {
		t.Errorf("Expected exit status message message \"exit status 0\" but got \"%s\"\n", e2.Error())
	}
}

func TestGetInput(t *testing.T) {
	input := "test"
	node := NewNode("localhost:0000")
	cli := NewCli(node, strings.NewReader(input + "\n"))
	returnedInput := cli.getInput()
	if returnedInput != input {
		t.Errorf("Expected input \"%s\" but got \"%s\"\n", input, returnedInput)
	}
}

func TestHandleInput(t *testing.T) {
	node := NewNode("localhost:0000")

	input1 := "exit"
	cli1 := NewCli(node, os.Stdin)
	if os.Getenv("BE_EXIT") == "1" {
		cli1.handleInput(input1, false)
		return
    }
    cmd1 := exec.Command(os.Args[0], "-test.run=TestExit")
    cmd1.Env = append(os.Environ(), "BE_EXIT=1")
    err1 := cmd1.Run()
	e1, exitedIncorrectly1 := err1.(*exec.ExitError)
	if exitedIncorrectly1 {
		t.Errorf("Expected exit status message message \"exit status 0\" but got \"%s\"\n", e1.Error())
	}

	input2 := "invalid input"
	cli2 := NewCli(node, os.Stdin)
	cli2.handleInput(input2, false)
}