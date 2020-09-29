package d7024e

import (
	"testing"
	"fmt"
)

func TestNewKademliaID(t *testing.T) {
	validID := leftPad("", '0', IDLength * 2)
	validKademliaID, err1 := NewKademliaID(validID)
	kIDType := fmt.Sprintf("%T", validKademliaID)
	if kIDType != "*d7024e.KademliaID" {
		t.Errorf("The returned object is not of type \"*d7024e.KademliaID\"")
	}
	if err1 != nil {
		t.Error(err1)
	}
	invalidID := leftPad("", '0', IDLength * 2 - 1)
	_, err2 := NewKademliaID(invalidID)
	if err2 == nil {
		t.Error("No error was returned given an invalid hexID")
	}
}

func TestNewRandomKademliaID(t *testing.T) {
	kID := NewRandomKademliaID()
	kIDType := fmt.Sprintf("%T", kID)
	if kIDType != "*d7024e.KademliaID" {
		t.Errorf("The returned object is not of type \"*d7024e.KademliaID\"")
	}
}

func TestLess(t *testing.T) {
	largeID := leftPad("", 'f', IDLength * 2)
	smallID := leftPad("", '0', IDLength * 2)
	largeKademliaID, _ := NewKademliaID(largeID)
	smallKademliaID, _ := NewKademliaID(smallID)
	if largeKademliaID.Less(smallKademliaID) {
		t.Errorf("The ID \"%s\" was calculated to be less than the ID \"%s\" but it is not", largeKademliaID, smallKademliaID)
	}
	if !smallKademliaID.Less(largeKademliaID) {
		t.Errorf("The ID \"%s\" was calculated to not be less than the ID \"%s\" but it is", smallKademliaID, largeKademliaID)
	}
	if largeKademliaID.Less(largeKademliaID) {
		t.Error("Comparing an ID to itself should return false but returned true")
	}
}

func TestEquals(t *testing.T) {
	largeID := leftPad("", 'f', IDLength * 2)
	smallID := leftPad("", '0', IDLength * 2)
	largeKademliaID, _ := NewKademliaID(largeID)
	smallKademliaID, _ := NewKademliaID(smallID)
	if largeKademliaID.Equals(smallKademliaID) {
		t.Errorf("The ID \"%s\" was calculated to be equal to the ID \"%s\" but it is not", largeKademliaID, smallKademliaID)
	}
	if smallKademliaID.Equals(largeKademliaID) {
		t.Errorf("The ID \"%s\" was calculated to be equal to the ID \"%s\" but it is not", smallKademliaID, largeKademliaID)
	}
	if !largeKademliaID.Equals(largeKademliaID) {
		t.Error("Comparing an ID to itself should return false but returned true")
	}
}

func TestEqualsZero(t *testing.T) {
	validID := leftPad("", '0', IDLength * 2)
	invalidID := leftPad("", 'a', IDLength * 2)
	validKademliaID, _ := NewKademliaID(validID)
	invalidKademliaID, _ := NewKademliaID(invalidID)
	if !validKademliaID.EqualsZero() {
		t.Errorf("The ID \"%s\" was expected to be equal to zero but was not", validKademliaID)
	}
	if invalidKademliaID.EqualsZero() {
		t.Errorf("The ID \"%s\" was expected to not be equal to zero but was", invalidKademliaID)
	}
}

func TestCalcDistance(t *testing.T) {
	largeID := leftPad("", 'a', IDLength * 2)
	smallID := leftPad("", '1', IDLength * 2)
	expectedID := leftPad("", 'b', IDLength * 2)
	largeKademliaID, _ := NewKademliaID(largeID)
	smallKademliaID, _ := NewKademliaID(smallID)
	expectedKademliaID, _ := NewKademliaID(expectedID)
	distance := largeKademliaID.CalcDistance(smallKademliaID)
	if !distance.Equals(expectedKademliaID) {
		t.Errorf("Distance between \"%s\" and \"%s\" should be \"%s\" but was \"%s\"", largeKademliaID, smallKademliaID, expectedKademliaID, distance)
	}
}

func TestString(t *testing.T) {
	strID := leftPad("", 'a', IDLength * 2)
	kademliaID, _ := NewKademliaID(strID)
	returnedStrID := kademliaID.String()
	if returnedStrID != strID {
		t.Errorf("Expected a string \"%s\" but got \"%s\"", strID, returnedStrID)
	}
}

func TestIDWithinRange(t *testing.T) {
	myID := leftPad("", '0', IDLength * 2)
	neighbourID := leftPad("1", '0', IDLength * 2)
	myKademliaID, _ := NewKademliaID(myID)
	neighbourKademliaID, _ := NewKademliaID(neighbourID)
	newKademliaID := myKademliaID.IDWithinRange()
	kIDType := fmt.Sprintf("%T", newKademliaID)
	if kIDType != "*d7024e.KademliaID" {
		t.Errorf("The returned object is not of type \"*d7024e.KademliaID\"")
	}
	if newKademliaID.Equals(myKademliaID)  {
		t.Errorf("The returned ID \"%s\" was equal to the object's KademliaID", newKademliaID.String())
	}
	if newKademliaID.Equals(neighbourKademliaID) {
		t.Errorf("The returned ID \"%s\" was equal to the object's neighbour's KademliaID", neighbourKademliaID.String())
	}
}