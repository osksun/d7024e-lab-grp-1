package d7024e

import (
	"testing"
)

func TestMatches(t *testing.T) {
	triggers := []string{"test", "t"}
	cmd := Cmd{
		triggers: triggers,
		description: "Test cmd for unit testing",
		usage: "\"test\", \"t\"",
		action: func(cli *Cli, args ...string) string {
			result := "Test cmd for unit testing result"
			return result
		},
	}
	if !cmd.matches(triggers[0]) {
		t.Errorf("Expected to find a match of the trigger \"%s\" when using triggers \"%s\" but match not found", triggers[0], triggers)
	}
	if !cmd.matches(triggers[1]) {
		t.Errorf("Expected to find a match of the trigger \"%s\" when using triggers \"%s\" but match not found", triggers[1], triggers)
	}
	invalidTrigger := "tes"
	if cmd.matches(invalidTrigger) {
		t.Errorf("Expected to not find a match of the trigger \"%s\" when using triggers \"%s\" but match found", invalidTrigger, triggers)
	}
}