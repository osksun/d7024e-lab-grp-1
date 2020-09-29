package d7024e

import (
	"testing"
)

func TestNewContact(t *testing.T){
	testKademliaID := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	testAddress := "localhost:8000"
	testContact := NewContact(testKademliaID, testAddress)

	if (testContact.Address != testAddress || testContact.ID != testKademliaID){
		t.Errorf("NewContact failed")
	}
}

func TestCalcDistance(t *testing.T){
	testKademliaID0 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	testKademliaID1 := NewKademliaID("FFFFFFFE00000000000000000000000000000000")
	testKademliaID2 := NewKademliaID("0000000010000000000000000000000000000000")
	testAddress := "localhost:8000"
	testContact := NewContact(testKademliaID0, testAddress)
	
	testContact.CalcDistance(testKademliaID1)
	if (testContact.Distance == testKademliaID2){
		t.Errorf("CalcDistance failed, distance was not 1")
	}


}

func TestLess(t *testing.T){
	bigID := leftPad("", "1", IDLength * 2)
}

