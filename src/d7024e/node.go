package d7024e

import (
	"net/http"
	"os"
)

// Node type definition
type Node struct {
	contact  *Contact
	rt       *RoutingTable
	vht      *ValueHashtable
	net      *Network
	kademlia *Kademlia
}


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
	node.kademlia = NewKademlia(node.rt, node.vht)
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
