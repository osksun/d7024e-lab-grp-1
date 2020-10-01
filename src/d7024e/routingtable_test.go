package d7024e

import (
	"fmt"
	"testing"
	"strconv"
)

func TestNewRoutingTable(t *testing.T) {
	kademliaID := NewRandomKademliaID()
	contact := NewContact(kademliaID, "localhost:0000")
	rt := NewRoutingTable(contact)
	rtType := fmt.Sprintf("%T", rt)
	if rtType != "*d7024e.RoutingTable" {
		t.Error("The returned object is not of expected type \"*d7024e.RoutingTable\"")
	}
	meType := fmt.Sprintf("%T", rt.me)
	if meType != "*d7024e.Contact" {
		t.Error("The returned object.me is not of expected type \"*d7024e.Contact\"")
	}
	bucketsType := fmt.Sprintf("%T", rt.buckets)
	bucketsExpectedType := "[" + strconv.Itoa(IDLength * 8) + "]*d7024e.bucket"
	if bucketsType != bucketsExpectedType {
		t.Errorf("The returned object.buckets is not of expected type \"%s\"\n", bucketsExpectedType)
	}
}

func TestRoutingTableAddContact(t *testing.T) {

}