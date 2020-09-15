package d7024e

import (
	"net/http"
	"time"
	"fmt"
)

type Node struct {
	contact *Contact
	rt *RoutingTable
	vht *ValueHashtable
	net *Network
	kademlia *Kademlia
}

const alpha = 1 // Alpha value should probably be stored in a Kademlia related file

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
	for {
		time.Sleep(1 * time.Second)
		if target != nil {
			fmt.Println("Contact found:", node.kademlia.LookupContact(target).Address)
		}		
	}
}

// AddContact adds an contact to the RoutingTable of the node
func (node *Node) AddContact(contact *Contact) {
	node.rt.AddContact(*contact)
}

// Contact returns the contact of the node
func (node *Node) Contact() *Contact {
	return node.contact
}
