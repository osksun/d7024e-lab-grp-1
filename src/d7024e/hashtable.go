// This code is taken from
// https://flaviocopes.com/golang-data-structure-hashtable/
// modified to use the static types string and []byte for key and value instead of the suggested generic types

// Package hashtable creates a ValueHashtable data structure for the Item type
package d7024e

import (
	"fmt"
	"sync"
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

// the hash() private function uses the famous Horner's method
// to generate a hash of a string with O(n) complexity
func hash(k string) int {
	key := fmt.Sprintf("%s", k)
	h := 0
	for i := 0; i < len(key); i++ {
		h = 31*h + int(key[i])
	}
	return h
}

// Put item with value v and key k into the hashtable
func (ht *ValueHashtable) Put(k string, v []byte) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	i := hash(k)
	if ht.items == nil {
		ht.items = make(map[int][]byte)
	}
	ht.items[i] = v
}

// Remove item with key k from hashtable
func (ht *ValueHashtable) Remove(k string) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	i := hash(k)
	delete(ht.items, i)
}

// Get item with key k from the hashtable
func (ht *ValueHashtable) Get(k string) []byte {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	i := hash(k)
	return ht.items[i]
}

// Size returns the number of the hashtable elements
func (ht *ValueHashtable) Size() int {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	return len(ht.items)
}
