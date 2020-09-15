package main

import (
    "./d7024e"
)

func main() {
	// Create nodes
	node1 := d7024e.NewNode("localhost:8000")
	node2 := d7024e.NewNode("localhost:8001")
	// Spinup nodes
	go node1.SpinupNode(node2.Contact, node2.Contact)
	node2.SpinupNode(nil, nil)
}