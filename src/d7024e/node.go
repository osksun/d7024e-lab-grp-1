package d7024e

import (
	"net/http"
	"time"
)

type Node struct {
	Contact *Contact
	Rt *RoutingTable
	Vht *ValueHashtable
	Net *Network
	Kademlia *Kademlia
}

// NewNode Constructor function for Node class
func NewNode(contact *Contact, alpha int) *Node {
	node := &Node{}
	node.Contact = contact
	node.Rt = NewRoutingTable(contact)
	node.Vht = NewValueHashtable()
	node.Net = NewNetwork(node.Rt, node.Vht)
	node.Kademlia = NewKademlia(node.Net, node.Rt, alpha)
	return node
}

// SpinupNode creates
func (node *Node) SpinupNode(target *Contact, receiver *Contact) {
	serveMux := http.NewServeMux()
	go node.Net.Listen(node.Contact.Address, serveMux)
	for {
		if target != nil {
			node.Net.SendFindContactMessage(target, receiver)
		}		
		time.Sleep(2 * time.Second)
	}
}