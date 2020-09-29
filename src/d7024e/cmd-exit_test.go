package d7024e

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestExit(t *testing.T) {
	node := NewNode("localhost:0000")
	cli := NewCli(node, strings.NewReader("exit\n"))
	if os.Getenv("BE_EXIT") == "1" {
		cli.Run(true, false)
		return
    }
    cmd := exec.Command(os.Args[0], "-test.run=TestExit")
    cmd.Env = append(os.Environ(), "BE_EXIT=1")
    err := cmd.Run()
	e, exitedIncorrectly := err.(*exec.ExitError)
	if exitedIncorrectly {
		t.Errorf("Expected exit status message message \"exit status 0\" but got \"%s\"\n", e.Error())
	}
}