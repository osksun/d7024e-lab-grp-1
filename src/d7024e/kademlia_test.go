package d7024e

import (
	"testing"
)

func TestNewKademlia(t *testing.T) {
	kademliaID := NewRandomKademliaID()
	contact := NewContact(kademliaID, "testAddress")
	rt := NewRoutingTable(contact)
	//vht := NewValueHashtable()
	//Net := NewNetwork(rt, vht)
	kademlia := NewKademlia(rt)
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
	node3.AddContact(node0.Contact())

	findNodeRequestChannel := make(chan findNodeRequest)
	go func() {
		for {
			RPCRequest := <- findNodeRequestChannel
			RPCResponse := findNodeResponse{}
			switch RPCRequest.receiver.Address {
			case node1.contact.Address:
				RPCResponse.contacts = node1.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node1.AddContact(node0.contact)
			case node2.contact.Address:
				RPCResponse.contacts = node2.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node2.AddContact(node0.contact)
			case node3.contact.Address:
				RPCResponse.contacts = node3.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node3.AddContact(node0.contact)
			case "":
				return
			}
			RPCResponse.sender = RPCRequest.receiver
			RPCRequest.responseChannel <- RPCResponse
		}
	}()

	test1TargetNode := node0.kademlia.LookupContact(node3.contact, findNodeRequestChannel)
	// Test to check whether it can find node3 starting from node0
	if (test1TargetNode.Address != node3.contact.Address){
		t.Error("LookupContact failed test 1, mismatch address")
	}

	test2TargetNode := node0.kademlia.LookupContact(node1.contact, findNodeRequestChannel)
	// Test to check whether node0 can find node1 in it's bucket and return it directly
	if (test2TargetNode.Address != node1.contact.Address){
		t.Error("Lookup failed test 2, mismatch address")
	}

	test3TargetNode := node0.kademlia.LookupContact(disconnectedNode.contact, findNodeRequestChannel)
	// Test to check whether node0 finds the closest node possible when looking for a node not in the network
	if (test3TargetNode.Address != node0.contact.Address){
		t.Error("Lookup failed test 3, mismatch address")
	}

	// Exit the goroutine
	findNodeRequestChannel <- findNodeRequest{NewContact(NewRandomKademliaID(), ""), NewContact(NewRandomKademliaID(), ""), nil}
}

func TestJoinNetwork(t *testing.T){
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
	node3.AddContact(node0.Contact())

	findNodeRequestChannel := make(chan findNodeRequest)
	go func() {
		for {
			RPCRequest := <- findNodeRequestChannel
			RPCResponse := findNodeResponse{}
			switch RPCRequest.receiver.Address {
			case node0.contact.Address:
				RPCResponse.contacts = node0.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node0.AddContact(disconnectedNode.contact)
				disconnectedNode.AddContact(node0.contact)
			case node1.contact.Address:
				RPCResponse.contacts = node1.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node1.AddContact(disconnectedNode.contact)
				disconnectedNode.AddContact(node1.contact)
			case node2.contact.Address:
				RPCResponse.contacts = node2.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node2.AddContact(disconnectedNode.contact)
				disconnectedNode.AddContact(node2.contact)
			case node3.contact.Address:
				RPCResponse.contacts = node3.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node3.AddContact(disconnectedNode.contact)
				disconnectedNode.AddContact(node3.contact)
			case "":
				return
			}
			RPCResponse.sender = RPCRequest.receiver
			RPCRequest.responseChannel <- RPCResponse
		}
	}()

	buckets := disconnectedNode.rt.buckets
	nContactsBefore := 0
	for _, bucket := range buckets {
		nContactsBefore += bucket.Len()
	}
	disconnectedNode.kademlia.JoinNetwork(node0.contact.Address, findNodeRequestChannel)
	nContactsAfter := 0
	for _, bucket := range buckets {
		nContactsAfter += bucket.Len()
	}
	if nContactsAfter <= nContactsBefore{
		t.Errorf("Expected node.rt.buckets to contain more contacts after joining the network but did not, went from %d to %d", nContactsBefore, nContactsAfter)
	}

	// Exit the goroutine
	findNodeRequestChannel <- findNodeRequest{NewContact(NewRandomKademliaID(), ""), NewContact(NewRandomKademliaID(), ""), nil}
}
