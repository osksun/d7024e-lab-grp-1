package d7024e
// This code is taken from
// https://flaviocopes.com/golang-data-structure-hashtable/
import (
	"sync"
	"crypto/sha1"
)

// HashSize the byte size of a hash
const HashSize = 20

// ValueHashtable the set of Items
type ValueHashtable struct {
	items map[[HashSize]byte][]byte
	lock  sync.RWMutex
}

// NewValueHashtable Constructor function for NewValueHashtable class
func NewValueHashtable() *ValueHashtable {
	valueHashtable := &ValueHashtable{}
	return valueHashtable
}

// Hash function uses sha1 hash function
func Hash(k []byte) [HashSize]byte {
	hasher := sha1.New()
	hasher.Write(k)
	hashSlice := hasher.Sum(nil)
	var hash [HashSize]byte;
	copy(hash[:], hashSlice)
	return hash
}


// Put item with value v and key k into the hashtable
func (ht *ValueHashtable) Put(k [HashSize]byte, v []byte) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	if ht.items == nil {
		ht.items = make(map[[HashSize]byte][]byte)
	}
	ht.items[k] = v
}

// Remove item with key k from hashtable
func (ht *ValueHashtable) Remove(k [HashSize]byte) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	delete(ht.items, k)
}

// Get item with key k from the hashtable
func (ht *ValueHashtable) Get(k [HashSize]byte) []byte {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	return ht.items[k]
}

// Size returns the number of the hashtable elements
func (ht *ValueHashtable) Size() int {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	return len(ht.items)
}
