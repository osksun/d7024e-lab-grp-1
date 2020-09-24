package d7024e

import (
	"encoding/hex"
	"math/rand"
)

// the static number of bytes in a KademliaID
const IDLength = 20

// type definition of a KademliaID
type KademliaID [IDLength]byte

// NewKademliaID returns a new instance of a KademliaID based on the string input
func NewKademliaID(data string) *KademliaID {
	decoded, _ := hex.DecodeString(data)

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	return &newKademliaID
}

// NewRandomKademliaID returns a new instance of a random KademliaID,
// change this to a better version if you like
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

// Less returns true if kademliaID < otherKademliaID (bitwise)
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if kademliaID == otherKademliaID (bitwise)
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// EqualsZero returns true if kademliaID is zero
func (kademliaID KademliaID) EqualsZero() bool {
	return kademliaID.Equals(NewKademliaID("0000000000000000000000000000000000000000"))
}

// CalcDistance returns a new instance of a KademliaID that is built
// through a bitwise XOR operation betweeen kademliaID and target
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

// String returns a simple string representation of a KademliaID
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}

// Returns an ID which is not itself or a neighbour
func (kademliaID *KademliaID) IDwithinRange() *KademliaID {
	var flag bool = true
	var resultKademliaID *KademliaID
	var neighbouringID *KademliaID = kademliaID
	var firstBitMask byte = 1

	if kademliaID[IDLength-1]&firstBitMask == 0 {
		neighbouringID[IDLength-1]++
	} else {
		neighbouringID[IDLength-1]--
	}
	for ok := true; ok; ok = flag {
		resultKademliaID = NewRandomKademliaID()
		if !(resultKademliaID.Equals(kademliaID)) || !(resultKademliaID.Equals(neighbouringID)) {
			flag = false
		}
	}
	return resultKademliaID
}
