package d7024e

import (
	"encoding/hex"
	"math/rand"
	"errors"
)

// IDLength the static number of bytes in a KademliaID
const IDLength = 20

// KademliaID type definition
type KademliaID [IDLength]byte

// NewKademliaID returns a new instance of a KademliaID based on the string input
func NewKademliaID(hexID string) (*KademliaID, error) {
	decoded, _ := hex.DecodeString(hexID)
	if len(decoded) != IDLength {
		return nil, errors.New("Invalid id length")
	}
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}
	return &newKademliaID, nil
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
	zID, _ := NewKademliaID("0000000000000000000000000000000000000000")
	return kademliaID.Equals(zID)
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

// NewKademliaIDWithinRange returns an ID which is not itself or a neighbour
func (kademliaID *KademliaID) NewKademliaIDWithinRange() *KademliaID {
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