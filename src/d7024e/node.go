package d7024e

import (
	"fmt"
	"net/http"
	"time"
)

type Node struct {
	contact  *Contact
	rt       *RoutingTable
	vht      *ValueHashtable
	net      *Network
	kademlia *Kademlia
}

const alpha = 3 // Alpha value should probably be stored in a Kademlia related file

// NewNode Constructor function for Node class
func NewNode(address string) *Node {
	node := &Node{}
	kademliaID := NewRandomKademliaID()
	node.contact = NewContact(kademliaID, address)
	node.rt = NewRoutingTable(node.contact)
	node.vht = NewValueHashtable()
	node.net = NewNetwork(node.rt, node.vht)
	node.kademlia = NewKademlia(node.net, node.rt, alpha)
	return node
}

// SpinupNode creates http server and listens on contact address
// Current parameters are temporary for testing to send messages
func (node *Node) SpinupNode(target *Contact) {
	serveMux := http.NewServeMux()
	go node.net.Listen(node.contact.Address, serveMux)
	time.Sleep(1 * time.Second)
	if target != nil {
		hash := node.kademlia.Store([]byte("ugabuga"), []byte("yo"))
		data := node.kademlia.LookupData(hash)
		fmt.Printf("Data found: \"%s\"", string(data))
		//NewCli(node).Run()
	}
	//time.Sleep(2 * time.Second)

}

// AddContact adds an contact to the RoutingTable of the node
func (node *Node) AddContact(contact *Contact) {
	node.rt.AddContact(*contact)
}

func (node *Node) JoinNetwork(address string) {
	// TODO Check if node is participating or not
	var refreshContact Contact // dummy contact only used for looking up random ID
	kademliaID := NewRandomKademliaID()
	node.contact = NewContact(kademliaID, node.contact.Address)
	//node.net.SendPingMessage(address, kademliaID)
	var nConnections = 0
	var nodeBuckets = node.rt.Buckets()
	node.kademlia.LookupContact(node.contact)
	for i := 0; i < len(nodeBuckets); i++ {
		nConnections += nodeBuckets[i].Len()
	}
	fmt.Println("Connections after 1st lookup: ", nConnections)
	refreshContact.ID = node.contact.ID.IDwithinRange()
	node.kademlia.LookupContact(&refreshContact)

	nConnections = 0
	for i := 0; i < len(nodeBuckets); i++ {
		nConnections += nodeBuckets[i].Len()
	}
	fmt.Println("Connections after 2nd lookup: ", nConnections)

}

// Contact returns the contact of the node
func (node *Node) Contact() *Contact {
	return node.contact
}

// Rt returns the routing table of the node
func (node *Node) Rt() *RoutingTable {
	return node.rt
}

// Vht returns the routing table of the node
func (node *Node) Vht() *ValueHashtable {
	return node.vht
}
