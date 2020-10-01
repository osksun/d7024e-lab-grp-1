package d7024e

import (
	"testing"
)

func TestLeftPad(t *testing.T) {
	str := "abc123"
	expectedResult := "0000" + str
	result := leftPad(str, '0', len(expectedResult))
	if result != expectedResult {
		t.Errorf("Got string \"%s\" but expected string \"%s\"", result, expectedResult)
	}
}
