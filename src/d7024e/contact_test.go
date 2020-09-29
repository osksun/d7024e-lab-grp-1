package d7024e

import (
	"testing"
)

func TestNewContact(t *testing.T){
	testKademliaID := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	testAddress := "localhost:8000"
	testContact := NewContact(testKademliaID, testAddress)

	if (testContact.Address != testAddress || testContact.ID != testKademliaID){
		t.Error("NewContact failed")
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
		t.Error("CalcDistance failed, distance was not 1")
	}


}

func TestLess(t *testing.T){
	bigID := leftPad("", 'F', IDLength * 2)
	mediumID := leftPad("", '8', IDLength * 2)
	tinyID := leftPad("", '0', IDLength * 2)
	bigContact := NewKademliaID(bigID)
	mediumContact := NewKademliaID(mediumID)
	tinyContact := NewKademliaID(tinyID)


	if tinyContact.Less(bigContact) == false{
		t.Error("Less failed, returned false for a distance larger than itself")
	}

	if mediumContact.Less(tinyContact) == true {
		t.Error("Less failed, returned true for a distance smaller than itself")
	}

	if bigContact.Less(bigContact) == true {
		t.Error("Less failed, returned true for a distance 0")
	}

}

func TestString(t *testing.T){
	testID := leftPad("", 'f', IDLength * 2)
	testKademliaID := NewKademliaID(testID)
	testAddress := "localhost:8000"
	testContact := NewContact(testKademliaID, testAddress)

	testString := testContact.String()
	if(testString != "contact(\"" + testID + "\", \"" + testAddress + "\")"){
		t.Error("String failed, string mismatch")
	}

}

func TestAppend(t *testing.T){
	testID0 := leftPad("", 'f', IDLength * 2)
	testKademliaID0 := NewKademliaID(testID0)
	testAddress0 := "localhost:8000"
	testContact0 := NewContact(testKademliaID0, testAddress0)

	testID1 := leftPad("", 'e', IDLength * 2)
	testKademliaID1 := NewKademliaID(testID1)
	testAddress1 := "localhost:8001"
	testContact1 := NewContact(testKademliaID1, testAddress1)

	var testContacts0 []Contact
	testContacts0 = append(testContacts0, *testContact0)
	testContacts0 = append(testContacts0, *testContact1)

	var testContactCandidates ContactCandidates
	testContactCandidates.Append(testContacts0)

	if (testContactCandidates.Len() != 2){
		t.Error("Append failed, incorrect size")
	}


}

func TestGetContacts(t *testing.T){
	testID0 := leftPad("", 'f', IDLength * 2)
	testKademliaID0 := NewKademliaID(testID0)
	testAddress0 := "localhost:8000"
	testContact0 := NewContact(testKademliaID0, testAddress0)

	testID1 := leftPad("", 'e', IDLength * 2)
	testKademliaID1 := NewKademliaID(testID1)
	testAddress1 := "localhost:8001"
	testContact1 := NewContact(testKademliaID1, testAddress1)

	var testContacts0 []Contact
	testContacts0 = append(testContacts0, *testContact0)
	testContacts0 = append(testContacts0, *testContact1)

	var testContactCandidates ContactCandidates
	testContactCandidates.Append(testContacts0)
	receivedContacts0 := testContactCandidates.GetContacts(1)
	if (len(receivedContacts0) != 1){
		t.Error("GetContacts failed, Incorrect size of contacts received")
	} 
	receivedContacts1 := testContactCandidates.GetContacts(2)
	if (len(receivedContacts1) != 2){
		t.Error("GetContacts failed, Incorrect size of contacts received")
	} 

}

func TestSort(t *testing.T){

	testID0 := leftPad("", 'f', IDLength * 2)
	testKademliaID0 := NewKademliaID(testID0)
	testAddress0 := "localhost:8000"
	testContact0 := NewContact(testKademliaID0, testAddress0)
	testContact0.CalcDistance(testContact0.ID)

	testID1 := leftPad("", '2', IDLength * 2)
	testKademliaID1 := NewKademliaID(testID1)
	testAddress1 := "localhost:8001"
	testContact1 := NewContact(testKademliaID1, testAddress1)
	testContact1.CalcDistance(testContact0.ID)

	testID2 := leftPad("", 'b', IDLength * 2)
	testKademliaID2 := NewKademliaID(testID2)
	testAddress2 := "localhost:8002"
	testContact2 := NewContact(testKademliaID2, testAddress2)
	testContact2.CalcDistance(testContact0.ID)

	testID3 := leftPad("", '5', IDLength * 2)
	testKademliaID3 := NewKademliaID(testID3)
	testAddress3 := "localhost:8003"
	testContact3 := NewContact(testKademliaID3, testAddress3)
	testContact3.CalcDistance(testContact0.ID)


	var testContacts0 []Contact
	testContacts0 = append(testContacts0, *testContact0)
	testContacts0 = append(testContacts0, *testContact1)
	testContacts0 = append(testContacts0, *testContact2)
	testContacts0 = append(testContacts0, *testContact3)


	var testContactCandidates ContactCandidates
	testContactCandidates.Append(testContacts0)
	testContactCandidates.Sort()
	for i := 0 ; i < testContactCandidates.Len() - 1; i++{
		if testContactCandidates.contacts[i + 1].Less(&testContactCandidates.contacts[i]) {
			t.Error("Sort failed, contacts not sorted")
		}
	}

	testContact0.CalcDistance(testContact0.ID)
	testContact1.CalcDistance(testContact1.ID)
	testContact2.CalcDistance(testContact2.ID)
	testContact3.CalcDistance(testContact3.ID)

	testContactCandidates.Sort()
	for i := 0 ; i < testContactCandidates.Len() - 1; i++{
		if testContactCandidates.contacts[i + 1].Less(&testContactCandidates.contacts[i]) {
			t.Error("Sort failed, contacts not sorted")
		}
	}


}

func TestContactLen(t *testing.T){
	testID0 := leftPad("", 'f', IDLength * 2)
	testKademliaID0 := NewKademliaID(testID0)
	testAddress0 := "localhost:8000"
	testContact0 := NewContact(testKademliaID0, testAddress0)
	testContact0.CalcDistance(testContact0.ID)

	testID1 := leftPad("", 'e', IDLength * 2)
	testKademliaID1 := NewKademliaID(testID1)
	testAddress1 := "localhost:8001"
	testContact1 := NewContact(testKademliaID1, testAddress1)
	testContact1.CalcDistance(testContact0.ID)

	testID2 := leftPad("", '1', IDLength * 2)
	testKademliaID2 := NewKademliaID(testID2)
	testAddress2 := "localhost:8002"
	testContact2 := NewContact(testKademliaID2, testAddress2)
	testContact2.CalcDistance(testContact0.ID)

	testID3 := leftPad("", '7', IDLength * 2)
	testKademliaID3 := NewKademliaID(testID3)
	testAddress3 := "localhost:8003"
	testContact3 := NewContact(testKademliaID3, testAddress3)
	testContact3.CalcDistance(testContact0.ID)


	var testContacts0 []Contact
	testContacts0 = append(testContacts0, *testContact0)
	testContacts0 = append(testContacts0, *testContact1)
	testContacts0 = append(testContacts0, *testContact2)
	testContacts0 = append(testContacts0, *testContact3)


	var testContactCandidates ContactCandidates
	if (len(testContactCandidates.contacts) != 0) {
		t.Error("Len failed, length should be 0")
	}
	testContactCandidates.Append(testContacts0)

	if (len(testContactCandidates.contacts) != 4) {
		t.Error("Len failed, length should be 4")
	}
}

func TestSwap(t *testing.T){
	testID0 := leftPad("", 'f', IDLength * 2)
	testKademliaID0 := NewKademliaID(testID0)
	testAddress0 := "localhost:8000"
	testContact0 := NewContact(testKademliaID0, testAddress0)
	testContact0.CalcDistance(testContact0.ID)

	testID1 := leftPad("", 'e', IDLength * 2)
	testKademliaID1 := NewKademliaID(testID1)
	testAddress1 := "localhost:8001"
	testContact1 := NewContact(testKademliaID1, testAddress1)
	testContact1.CalcDistance(testContact0.ID)

	testID2 := leftPad("", '1', IDLength * 2)
	testKademliaID2 := NewKademliaID(testID2)
	testAddress2 := "localhost:8002"
	testContact2 := NewContact(testKademliaID2, testAddress2)
	testContact2.CalcDistance(testContact0.ID)

	testID3 := leftPad("", '2', IDLength * 2)
	testKademliaID3 := NewKademliaID(testID3)
	testAddress3 := "localhost:8003"
	testContact3 := NewContact(testKademliaID3, testAddress3)
	testContact3.CalcDistance(testContact0.ID)

	var testContacts0 []Contact
	testContacts0 = append(testContacts0, *testContact0)
	testContacts0 = append(testContacts0, *testContact1)
	testContacts0 = append(testContacts0, *testContact2)
	testContacts0 = append(testContacts0, *testContact3)

	var testContactCandidates ContactCandidates
	testContactCandidates.Append(testContacts0)
	testContactCandidates.Swap(0,3)

	if (testContactCandidates.contacts[0].Address != "localhost:8003"){
		t.Error("Swap failed, address mismatch")
	}
}

func TestContactCandidatesLess (t *testing.T){
	testID0 := leftPad("", 'f', IDLength * 2)
	testKademliaID0 := NewKademliaID(testID0)
	testAddress0 := "localhost:8000"
	testContact0 := NewContact(testKademliaID0, testAddress0)
	testContact0.CalcDistance(testContact0.ID)

	testID1 := leftPad("", 'e', IDLength * 2)
	testKademliaID1 := NewKademliaID(testID1)
	testAddress1 := "localhost:8001"
	testContact1 := NewContact(testKademliaID1, testAddress1)
	testContact1.CalcDistance(testContact0.ID)

	if (testContact1.Less(testContact0) == true ){
		t.Error("Less failed, returned true for distance greater than target ")
	}
}

func TestRemoveDuplicates(t *testing.T){
	testID0 := leftPad("", 'f', IDLength * 2)
	testKademliaID0 := NewKademliaID(testID0)
	testAddress0 := "localhost:8000"
	testContact0 := NewContact(testKademliaID0, testAddress0)

	testID1 := leftPad("", 'e', IDLength * 2)
	testKademliaID1 := NewKademliaID(testID1)
	testAddress1 := "localhost:8000"
	testContact1 := NewContact(testKademliaID1, testAddress1)

	testID2 := leftPad("", '6', IDLength * 2)
	testKademliaID2 := NewKademliaID(testID2)
	testAddress2 := "localhost:8000"
	testContact2 := NewContact(testKademliaID2, testAddress2)

	testID3 := leftPad("", '1', IDLength * 2)
	testKademliaID3 := NewKademliaID(testID3)
	testAddress3 := "localhost:8000"
	testContact3 := NewContact(testKademliaID3, testAddress3)

	testID4 := leftPad("", '2', IDLength * 2)
	testKademliaID4 := NewKademliaID(testID4)
	testAddress4 := "localhost:8003"
	testContact4 := NewContact(testKademliaID4, testAddress4)

	var testContacts0 []Contact
	testContacts0 = append(testContacts0, *testContact0)
	testContacts0 = append(testContacts0, *testContact1)
	testContacts0 = append(testContacts0, *testContact2)
	testContacts0 = append(testContacts0, *testContact3)
	testContacts0 = append(testContacts0, *testContact4)

	var testContactCandidates ContactCandidates
	testContactCandidates.Append(testContacts0)
	testContactCandidates.RemoveDuplicates()
	if (testContactCandidates.Len() != 2) {
		t.Error("RemoveDuplicated failed, incorrect size after operation")
	}
}
