// This code is taken from
// https://flaviocopes.com/golang-data-structure-hashtable/
// modified to use the static types string and []byte for key and value instead of the suggested generic types

// Package hashtable creates a ValueHashtable data structure for the Item type
package d7024e

import (
	"sync"
	"crypto/sha1"
	//	"github.com/cheekybits/genny/generic"
)

// Key the key of the dictionary
//type Key string

// Value the content of the dictionary
//type Value []byte

// ValueHashtable the set of Items
type ValueHashtable struct {
	items map[int][]byte
	lock  sync.RWMutex
}

// NewValueHashtable Constructor function for NewValueHashtable class
func NewValueHashtable() *ValueHashtable {
	valueHashtable := &ValueHashtable{}
	return valueHashtable
}

// the hash() private function uses sha1 hash function
func hash(k []byte) [20]byte {
	hasher := sha1.New()
	hasher.Write(k)
	hashSlice := hasher.Sum(nil)
	var hash [20]byte;
	copy(hash[:], hashSlice)
	return hash
}


// Put item with value v and key k into the hashtable
func (ht *ValueHashtable) Put(k []byte, v []byte) [20]byte {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	hash := hash(k)
	if ht.items == nil {
		ht.items = make(map[[20]byte][]byte)
	}
	ht.items[hash] = v
	return hash
}

// Remove item with key k from hashtable
func (ht *ValueHashtable) Remove(k []byte) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	hash := hash(k)
	delete(ht.items, hash)
}

// Get item with key k from the hashtable
func (ht *ValueHashtable) Get(k []byte) []byte {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	hash := hash(k)
	return ht.items[hash]
}

// Size returns the number of the hashtable elements
func (ht *ValueHashtable) Size() int {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	return len(ht.items)
}
