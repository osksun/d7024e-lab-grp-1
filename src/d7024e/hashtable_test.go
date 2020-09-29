// // This code is taken from
// // https://flaviocopes.com/golang-data-structure-hashtable/
// // modified to use the static types string and []byte for key and value instead of the suggested generic types

// // Package hashtable creates a ValueHashtable data structure for the Item type
package d7024e

import (
	"fmt"
	"testing"
)

func TestNewValueHashtable(t *testing.T) {
	ht := NewValueHashtable()
	htType := fmt.Sprintf("%T", ht)
	if htType != "*d7024e.ValueHashtable" {
		t.Errorf("The hashtable is not of type ValueHashtable")
	}
	fmt.Printf("TestNewValueHashtable finished running with status OK\n")
}

func TestHash(t *testing.T) {
	k1 := []byte("key")
	k2 := []byte("key")
	h1 := Hash(k1)
	h2 := Hash(k2)
	if h1 != h2 {
		t.Errorf("Hash produced is not consistent with input.")
	}
	fmt.Printf("TestHash finished running with status OK\n")
}

func TestPut(t *testing.T) {
	var k [HashSize]byte
	k1 := []byte("test")
	v1 := []byte("value")
	ht := NewValueHashtable()
	copy(k[:], k1)
	ht.Put(k, v1)
	fmt.Println(ht.items)
}

// // Put item with value v and key k into the hashtable
// func (ht *ValueHashtable) Put(k [HashSize]byte, v []byte) {
// 	ht.lock.Lock()
// 	defer ht.lock.Unlock()
// 	if ht.items == nil {
// 		ht.items = make(map[[HashSize]byte][]byte)
// 	}
// 	ht.items[k] = v
// }

// // Remove item with key k from hashtable
// func (ht *ValueHashtable) Remove(k [HashSize]byte) {
// 	ht.lock.Lock()
// 	defer ht.lock.Unlock()
// 	delete(ht.items, k)
// }

// // Get item with key k from the hashtable
// func (ht *ValueHashtable) Get(k [HashSize]byte) []byte {
// 	ht.lock.RLock()
// 	defer ht.lock.RUnlock()
// 	return ht.items[k]
// }

// // Size returns the number of the hashtable elements
// func (ht *ValueHashtable) Size() int {
// 	ht.lock.RLock()
// 	defer ht.lock.RUnlock()
// 	return len(ht.items)
// }
