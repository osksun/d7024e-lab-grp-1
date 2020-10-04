package d7024e

import (
	"testing"
	"reflect"
)

func TestNewKademlia(t *testing.T) {
	kademliaID := NewRandomKademliaID()
	contact := NewContact(kademliaID, "testAddress")
	rt := NewRoutingTable(contact)
	ht := NewValueHashtable()
	//Net := NewNetwork(rt, vht)
	kademlia := NewKademlia(rt, ht)
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

	node0.rt.AddContact(*node1.contact)
	node1.rt.AddContact(*node2.contact)
	node2.rt.AddContact(*node3.contact)
	node3.rt.AddContact(*node0.contact)

	findNodeRequestChannel := make(chan findNodeRequest)
	go func() {
		for {
			RPCRequest := <- findNodeRequestChannel
			RPCResponse := findNodeResponse{}
			switch RPCRequest.receiver.Address {
			case node1.contact.Address:
				RPCResponse.contacts = node1.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node1.rt.AddContact(*node0.contact)
			case node2.contact.Address:
				RPCResponse.contacts = node2.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node2.rt.AddContact(*node0.contact)
			case node3.contact.Address:
				RPCResponse.contacts = node3.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node3.rt.AddContact(*node0.contact)
			case "":
				return
			}
			RPCResponse.sender = *RPCRequest.receiver
			RPCRequest.responseChannel <- RPCResponse
		}
	}()

	test1TargetNode := node0.kademlia.LookupContact(node3.contact, 1, findNodeRequestChannel)[0]
	// Test to check whether it can find node3 starting from node0
	if (test1TargetNode.Address != node3.contact.Address){
		t.Error("LookupContact failed test 1, mismatch address")
	}

	test2TargetNode := node0.kademlia.LookupContact(node1.contact, 1, findNodeRequestChannel)[0]
	// Test to check whether node0 can find node1 in it's bucket and return it directly
	if (test2TargetNode.Address != node1.contact.Address){
		t.Error("Lookup failed test 2, mismatch address")
	}

	test3TargetNode := node0.kademlia.LookupContact(disconnectedNode.contact, 1, findNodeRequestChannel)[0]
	// Test to check whether node0 finds the closest node possible when looking for a node not in the network
	if (test3TargetNode.Address != node0.contact.Address){
		t.Error("Lookup failed test 3, mismatch address")
	}

	// Exit the goroutine
	findNodeRequestChannel <- findNodeRequest{NewContact(NewRandomKademliaID(), ""), NewContact(NewRandomKademliaID(), ""), nil}
}

func TestJoinNetwork(t *testing.T) {
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

	node0.rt.AddContact(*node1.contact)
	node1.rt.AddContact(*node2.contact)
	node2.rt.AddContact(*node3.contact)
	node3.rt.AddContact(*node0.contact)

	findNodeRequestChannel := make(chan findNodeRequest)
	go func() {
		for {
			RPCRequest := <- findNodeRequestChannel
			RPCResponse := findNodeResponse{}
			switch RPCRequest.receiver.Address {
			case node0.contact.Address:
				RPCResponse.contacts = node0.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node0.rt.AddContact(*disconnectedNode.contact)
				disconnectedNode.rt.AddContact(*node0.contact)
			case node1.contact.Address:
				RPCResponse.contacts = node1.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node1.rt.AddContact(*disconnectedNode.contact)
				disconnectedNode.rt.AddContact(*node1.contact)
			case node2.contact.Address:
				RPCResponse.contacts = node2.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node2.rt.AddContact(*disconnectedNode.contact)
				disconnectedNode.rt.AddContact(*node2.contact)
			case node3.contact.Address:
				RPCResponse.contacts = node3.rt.FindClosestContacts(RPCRequest.target.ID, k)
				node3.rt.AddContact(*disconnectedNode.contact)
				disconnectedNode.rt.AddContact(*node3.contact)
			case "":
				return
			}
			RPCResponse.sender = *RPCRequest.receiver
			RPCRequest.responseChannel <- RPCResponse
		}
	}()

	buckets := disconnectedNode.rt.buckets
	nContactsBefore := 0
	for i := 0; i < len(buckets); i++ {
		nContactsBefore += buckets[i].Len()
	}
	disconnectedNode.kademlia.JoinNetwork(node0.contact.Address, findNodeRequestChannel)
	nContactsAfter := 0
	for i := 0; i < len(buckets); i++ {
		nContactsAfter += buckets[i].Len()
	}
	if nContactsAfter <= nContactsBefore{
		t.Errorf("Expected node.rt.buckets to contain more contacts after joining the network but did not, went from %d to %d", nContactsBefore, nContactsAfter)
	}

	// Exit the goroutine
	findNodeRequestChannel <- findNodeRequest{NewContact(NewRandomKademliaID(), ""), NewContact(NewRandomKademliaID(), ""), nil}
}

func TestStore(t *testing.T) {
	kademliaID0 := leftPad("0", '0', IDLength * 2)
	kademliaID1 := leftPad("1", '0', IDLength * 2)
	kademliaID2 := leftPad("2", '0', IDLength * 2)
	kademliaID3 := leftPad("3", '0', IDLength * 2)
	kademliaID4 := leftPad("4", '0', IDLength * 2)
	node0 := NewNode("localhost:5000", kademliaID0)
	node1 := NewNode("localhost:5001", kademliaID1)
	node2 := NewNode("localhost:5002", kademliaID2)
	node3 := NewNode("localhost:5003", kademliaID3)
	unresponsiveNode := NewNode("localhost:5004", kademliaID4)

	node0.rt.AddContact(*node1.contact)
	node0.rt.AddContact(*unresponsiveNode.contact)
	node1.rt.AddContact(*node2.contact)
	node2.rt.AddContact(*node3.contact)
	node3.rt.AddContact(*node0.contact)

	findNodeRequestChannel := make(chan findNodeRequest)
	go func() {
		for {
			findNodeRequest := <- findNodeRequestChannel
			findNodeReponse := findNodeResponse{}
			switch findNodeRequest.receiver.Address {
			case node1.contact.Address:
				findNodeReponse.contacts = node1.rt.FindClosestContacts(findNodeRequest.target.ID, k)
				node1.rt.AddContact(*node0.contact)
				node0.rt.AddContact(*node1.contact)
			case node2.contact.Address:
				findNodeReponse.contacts = node2.rt.FindClosestContacts(findNodeRequest.target.ID, k)
				node2.rt.AddContact(*node0.contact)
				node0.rt.AddContact(*node2.contact)
			case node3.contact.Address:
				findNodeReponse.contacts = node3.rt.FindClosestContacts(findNodeRequest.target.ID, k)
				node3.rt.AddContact(*node0.contact)
				node0.rt.AddContact(*node3.contact)
			case "":
				return
			}
			findNodeReponse.sender = *findNodeRequest.receiver
			findNodeRequest.responseChannel <- findNodeReponse
		}
	}()

	storeRequestChannel := make(chan storeRequest)
	go func() {
		for {
			storeRequest := <- storeRequestChannel
			switch storeRequest.receiver.Address {
			case node0.contact.Address:
				node0.vht.Put(storeRequest.hash, storeRequest.data)
			case node1.contact.Address:
				node0.vht.Put(storeRequest.hash, storeRequest.data)
			case node2.contact.Address:
				node0.vht.Put(storeRequest.hash, storeRequest.data)
			case node3.contact.Address:
				node0.vht.Put(storeRequest.hash, storeRequest.data)
			case "":
				return
			}
		}
	}()

	storedCopiesBefore := node0.vht.Size() + node1.vht.Size() + node2.vht.Size() + node3.vht.Size()
	filename := []byte("filename")
	data := []byte("data")
	hash := node0.kademlia.Store(filename, data, storeRequestChannel, findNodeRequestChannel)
	storedCopiesAfter := node0.vht.Size() + node1.vht.Size() + node2.vht.Size() + node3.vht.Size()
	if storedCopiesAfter < storedCopiesBefore {
		t.Error("The total amount of stored copies of files on the network was not increased when performing kademlia.Store")
	}
	// Exit the goroutines
	findNodeRequestChannel <- findNodeRequest{NewContact(NewRandomKademliaID(), ""), NewContact(NewRandomKademliaID(), ""), nil}
	storeRequestChannel <- storeRequest{NewContact(NewRandomKademliaID(), ""), hash, data}
}

func TestLookupData(t *testing.T) {
	kademliaID0 := leftPad("0", '0', IDLength * 2)
	kademliaID1 := leftPad("1", '0', IDLength * 2)
	kademliaID2 := leftPad("2", '0', IDLength * 2)
	kademliaID3 := leftPad("3", '0', IDLength * 2)
	kademliaID4 := leftPad("4", '0', IDLength * 2)
	node0 := NewNode("localhost:5000", kademliaID0)
	node1 := NewNode("localhost:5001", kademliaID1)
	node2 := NewNode("localhost:5002", kademliaID2)
	node3 := NewNode("localhost:5003", kademliaID3)
	unresponsiveNode := NewNode("localhost:5004", kademliaID4)

	node0.rt.AddContact(*node1.contact)
	node0.rt.AddContact(*unresponsiveNode.contact)
	node1.rt.AddContact(*node2.contact)
	node2.rt.AddContact(*node3.contact)
	node3.rt.AddContact(*node0.contact)

	findDataRequestChannel := make(chan findDataRequest)
	go func() {
		for {
			findDataRequest := <- findDataRequestChannel
			findDataReponse := findDataResponse{}
			switch findDataRequest.receiver.Address {
			case node1.contact.Address:
				findDataReponse.data = node1.vht.Get(findDataRequest.hash)
				if findDataReponse.data == nil {
					findDataReponse.contacts = node1.rt.FindClosestContacts(findDataRequest.target.ID, k)
				}
				node1.rt.AddContact(*node0.contact)
				node0.rt.AddContact(*node1.contact)
			case node2.contact.Address:
				findDataReponse.data = node2.vht.Get(findDataRequest.hash)
				if findDataReponse.data == nil {
					findDataReponse.contacts = node2.rt.FindClosestContacts(findDataRequest.target.ID, k)
				}
				node2.rt.AddContact(*node0.contact)
				node0.rt.AddContact(*node2.contact)
			case node3.contact.Address:
				findDataReponse.data = node3.vht.Get(findDataRequest.hash)
				if findDataReponse.data == nil {
					findDataReponse.contacts = node3.rt.FindClosestContacts(findDataRequest.target.ID, k)
				}
				node3.rt.AddContact(*node0.contact)
				node0.rt.AddContact(*node3.contact)
			case unresponsiveNode.contact.Address:
				continue
			case "":
				return
			}
			findDataReponse.sender = *findDataRequest.receiver
			findDataRequest.responseChannel <- findDataReponse
		}
	}()
	filename := []byte("filename")
	data := []byte("data")
	hash := Hash(filename)
	node2.vht.Put(hash, data)
	dataRetrieved := node0.kademlia.LookupData(hash, findDataRequestChannel)
	if !reflect.DeepEqual(dataRetrieved, data) {
		t.Errorf("The retrieved data \"%s\" was not equal to expected data \"%s\"", string(dataRetrieved), string(data))
	}
	nonExistingFilename := []byte("nonExistingFilename")
	nonExistingHash := Hash(nonExistingFilename)
	badData := node0.kademlia.LookupData(nonExistingHash, findDataRequestChannel)
	if badData != nil {
		t.Errorf("The retrieved data \"%s\" was when looking up non existing hash", string(badData))
	}
	// Exit the goroutine
	findDataRequestChannel <- findDataRequest{hash, NewContact(NewRandomKademliaID(), ""), NewContact(NewRandomKademliaID(), ""), nil}
}