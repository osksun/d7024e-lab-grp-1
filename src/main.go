package main

import (
	"./d7024e"
	"strconv"
	"math/rand"
	"fmt"
)

func leftPad(str string, pad rune, lenght int) string {
	for i := len(str); i < lenght; i++ {
		str = string(pad) + str
	}
	return str
}

func main() {
	// Create nodes
	const nNodes = 50
	var nodes [nNodes]*d7024e.Node
	for i := 0; i < nNodes; i++ {
		nodes[i] = d7024e.NewNode("localhost:" + strconv.Itoa(i+10000))
	}
	// Create random connections
	rand.Seed(0)
	nConnections := 10000 // Note that this is not the same number as the number of final connnections due to not avoiding collisions
	for i := 0; i < nConnections; i++ {
		n := rand.Intn(nNodes)
		c := rand.Intn(nNodes - 1) 
		if (c >= n) {
			c++
		}
		nodes[n].AddContact(nodes[c].Contact())
	}
	
	// Count number of contacts in each node's bucktes
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

	targetContact := nodes[20].Contact()
	nodes[0].SpinupNode(targetContact)
}