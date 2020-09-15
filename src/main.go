package main

import (
    "./d7024e"
)

func main() {
	// Create nodes
	node1 := d7024e.NewNode("localhost:8001")
	node2 := d7024e.NewNode("localhost:8002")
	node3 := d7024e.NewNode("localhost:8003")
	// Create connections
	node1.AddContact(node2.Contact())
	node2.AddContact(node1.Contact())
	node2.AddContact(node3.Contact())
	// Spinup nodes
	go node1.SpinupNode(node3.Contact(), node2.Contact()) // Target, Receiver
	go node2.SpinupNode(nil, nil)
	go node3.SpinupNode(nil, nil)
	// Infinite loop to prevent program from exiting
	for {}
}