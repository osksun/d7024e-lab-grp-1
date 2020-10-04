package d7024e

import (
	"testing"
	"fmt"
)

// NewNode Constructor function for Node class
func TestNewNode(t *testing.T) {
	node := NewNode("localhost:8000", "")
	nodeType := fmt.Sprintf("%T", node)
	if nodeType != "*d7024e.Node" {
		t.Error("The node is not of type \"*d7024e.Node\"")
	}
	contactType := fmt.Sprintf("%T", node.contact)
	if contactType != "*d7024e.Contact" {
		t.Error("The node.contact is not of type \"*d7024e.Contact\"")
	}
	rtType := fmt.Sprintf("%T", node.rt)
	if rtType != "*d7024e.RoutingTable" {
		t.Error("The node.rt is not of type \"*d7024e.RoutingTable\"")
	}
	vhtType := fmt.Sprintf("%T", node.vht)
	if vhtType != "*d7024e.ValueHashtable" {
		t.Error("The node.vht is not of type \"*d7024e.ValueHashtable\"")
	}
	netType := fmt.Sprintf("%T", node.net)
	if netType != "*d7024e.Network" {
		t.Error("The node.net is not of type \"*d7024e.Network\"")
	}
	kademliaType := fmt.Sprintf("%T", node.kademlia)
	if kademliaType != "*d7024e.Kademlia" {
		t.Error("The node.kademlia is not of type \"*d7024e.Kademlia\"")
	}
}

func TestSpinupNode(t *testing.T) {

}