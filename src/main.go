package main

import (
	"./d7024e"
)

func main() {
	node := d7024e.NewNode("localhost:8000", "")
	node.SpinupNode(false, true)
}
