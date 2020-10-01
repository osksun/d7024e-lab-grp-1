package d7024e

import (
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	node := NewNode("localhost:0000")
	cli := NewCli(node, strings.NewReader("help\n"))
	statusCode := cli.Run(true, false)
	if statusCode != 0 {
		t.Errorf("Expected return with status code \"0\" but got \"%d\"", statusCode)
	}
}