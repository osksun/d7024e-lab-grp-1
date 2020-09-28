package main

import (
	/*"math/rand"
	"strconv"*/
	"fmt"
	"./d7024e"
)

func leftPad(str string, pad rune, lenght int) string {
	for i := len(str); i < lenght; i++ {
		str = string(pad) + str
	}
	return str
}

func main() {
	// Create nodes
	
	node1 := d7024e.NewNode("localhost:8000")
	node2 := d7024e.NewNode("localhost:8001")
	node3 := d7024e.NewNode("localhost:8002")
	node4 := d7024e.NewNode("localhost:8003")
	node5 := d7024e.NewNode("localhost:8004")

	node1.AddContact(node2.Contact())
	node2.AddContact(node3.Contact())
	node4.AddContact(node1.Contact())
	node3.AddContact(node5.Contact())
	node5.AddContact(node3.Contact())

	go node2.SpinupNode(nil)
	go node3.SpinupNode(nil)
	go node4.SpinupNode(nil)
	go node5.SpinupNode(nil)
	node1.SpinupNode(nil)
	
	
	var n4Buckets = node4.Rt().Buckets()
	nConnections := 0
	for i := 0; i < len(n4Buckets); i++ {
		nConnections += n4Buckets[i].Len()
	}
	fmt.Println("node4 connections before: ", nConnections)
	// Fake bucket insertion
	node1.AddContact(node2.Contact())
	node4.JoinNetwork("localhost:8000")

	nConnections = 0
	for i := 0; i < len(n4Buckets); i++ {
		nConnections += n4Buckets[i].Len()
	}
	fmt.Println("node4 connections after: ", nConnections)

	/*const nNodes = 500
	var nodes [nNodes]*d7024e.Node
	for i := 0; i < nNodes; i++ {
		nodes[i] = d7024e.NewNode("localhost:" + strconv.Itoa(i+10000))
	}

	// Create random connections
	rand.Seed(0)
	nConnections := 30000 // Note that this is not the same number as the number of final connnections due to not avoiding collisions
	for i := 0; i < nConnections; i++ {
		n := rand.Intn(nNodes)
		c := rand.Intn(nNodes - 1)
		if c >= n {
			c++
		}
		nodes[n].AddContact(nodes[c].Contact())
	}
	// Count number of contacts in each node's buckets
	for i := 0; i < nNodes; i++ {
		nConnections := 0
		buckets := nodes[i].Rt().Buckets()
		for j := 0; j < len(buckets); j++ {
			nConnections += buckets[j].Len()
		}
		fmt.Println(i, "has", nConnections, "connections")
	}

	// Spinup nodes
	for i := 1; i < nNodes; i++ {
		go nodes[i].SpinupNode(nil)
	}

	targetContact := nodes[40].Contact()
	nodes[0].SpinupNode(targetContact)*/
}