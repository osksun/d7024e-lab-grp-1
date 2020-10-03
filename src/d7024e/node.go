package d7024e

import (
	"net/http"
	"os"
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
func NewNode(address string, kademliaID string) *Node {
	node := &Node{}
	var kID *KademliaID
	if kademliaID != "" {
		kID, _ = NewKademliaID(kademliaID)
	} else {
		kID = NewRandomKademliaID()
	}
	node.contact = NewContact(kID, address)
	node.rt = NewRoutingTable(node.contact)
	node.vht = NewValueHashtable()
	node.net = NewNetwork(node.rt, node.vht)
	node.kademlia = NewKademlia(node.net, node.rt, alpha)
	return node
}

// SpinupNode creates http server and listens on contact address
func (node *Node) SpinupNode(cliRunOnce bool, cliVerbose bool) {
	serveMux := http.NewServeMux()
	go node.net.Listen(node.contact.Address, serveMux)
	go node.net.handleChannels()
	NewCli(node, os.Stdin).Run(cliRunOnce, cliVerbose)
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
	refreshContact.ID = node.contact.ID.IDWithinRange()
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
