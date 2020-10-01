package d7024e

import (
	"testing"
	"net/http"
)

func TestNewKademlia(t *testing.T) {
	kademliaID := NewRandomKademliaID()
	contact := NewContact(kademliaID, "testAddress")
	rt := NewRoutingTable(contact)
	vht := NewValueHashtable()
	Net := NewNetwork(rt, vht)

	kademlia := NewKademlia(Net , rt, 1)
	if (kademlia == nil){
		t.Errorf("NewKademlia failed, variable was nil")
	}
}

func TestLookupContact(t *testing.T){
	kademliaID0 := leftPad("0", '0', IDLength * 2)
	kademliaID1 := leftPad("1", '0', IDLength * 2)
	kademliaID2 := leftPad("2", '0', IDLength * 2)
	kademliaID3 := leftPad("3", '0', IDLength * 2)
	kademliaID4 := leftPad("4", '0', IDLength * 2)

	node0 := NewNode("localhost:5000", kademliaID0)
	node1 := NewNode("localhost:5001", kademliaID1)
	node2 := NewNode("localhost:5002", kademliaID2)
	node3 := NewNode("localhost:5003", kademliaID3)
	disconnectedNode := NewNode("localhost:5010", kademliaID4)

	node0.AddContact(node1.Contact())
	node1.AddContact(node2.Contact())
	node2.AddContact(node3.Contact())
	node3.AddContact(node1.Contact())

	serveMux0 := http.NewServeMux()
	serveMux1 := http.NewServeMux()
	serveMux2 := http.NewServeMux()
	serveMux3 := http.NewServeMux()
	go node0.net.Listen(node0.contact.Address, serveMux0)
	go node1.net.Listen(node1.contact.Address, serveMux1)
	go node2.net.Listen(node2.contact.Address, serveMux2)
	go node3.net.Listen(node3.contact.Address, serveMux3)

	test1TargetNode := node0.kademlia.LookupContact(node3.contact)
	// Test to check whether it can find node3 starting from node0
	if (test1TargetNode.Address != node3.contact.Address){
		t.Error("LookupContact failed test 1, mismatch address")
	}

	test2TargetNode := node0.kademlia.LookupContact(node1.contact)
	// Test to check whether node0 can find node1 in it's bucket and return it directly
	if (test2TargetNode.Address != node1.contact.Address){
		t.Error("Lookup failed test 2, mismatch address")
	}

	test3TargetNode := node0.kademlia.LookupContact(disconnectedNode.contact)
	// Test to check whether node0 finds the closest node possible when looking for a node not in the network 
	if (test3TargetNode.Address != node0.contact.Address){
		t.Error("Lookup failed test 3, mismatch address")
	}
}

func TestLookupData(t *testing.T){
	
}


