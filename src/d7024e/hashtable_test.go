// // This code is taken from
// // https://flaviocopes.com/golang-data-structure-hashtable/
// // modified to use the static types string and []byte for key and value instead of the suggested generic types

// // Package hashtable creates a ValueHashtable data structure for the Item type
package d7024e

import (
	"fmt"
	"reflect"
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

	it := ht.items

	copy(k[:], k1)
	ht.Put(k, v1)

	it2 := ht.items

	if it != nil || it2 == nil {
		t.Errorf("Hashtable items wasn't initialized as nil or not made.")
	}

	testitems := make(map[[HashSize]byte][]byte)
	testitems[k] = v1

	if !reflect.DeepEqual(ht.items, testitems) {
		t.Errorf("Hashtable didn't put correctly.")
	}
	fmt.Printf("TestPut finished running with status OK\n")
}

func TestRemove(t *testing.T) {
	var k [HashSize]byte
	k1 := []byte("test")
	v1 := []byte("value")
	ht := NewValueHashtable()
	copy(k[:], k1)

	testitems := make(map[[HashSize]byte][]byte)
	testitems[k] = v1

	ht.items = testitems
	ht.Remove(k)

	if ht.items[k] != nil {
		t.Errorf("Hashtable didn't remove correctly.")
	}
	fmt.Printf("TestRemove finished running with status OK\n")
}

func TestGet(t *testing.T) {
	var k [HashSize]byte
	k1 := []byte("test")
	v1 := []byte("value")
	ht := NewValueHashtable()
	copy(k[:], k1)

	testitems := make(map[[HashSize]byte][]byte)
	testitems[k] = v1

	ht.items = testitems

	item := ht.Get(k)
	if !reflect.DeepEqual(item, v1) {
		t.Errorf("Hashtable didn't get correctly.")
	}
	fmt.Printf("TestGet finished running with status OK\n")
}

func TestSize(t *testing.T) {
	var k [HashSize]byte
	k1 := []byte("test")
	v1 := []byte("value")
	ht := NewValueHashtable()
	copy(k[:], k1)

	testitems := make(map[[HashSize]byte][]byte)
	testitems[k] = v1

	ht.items = testitems
	if ht.Size() != 1 {
		t.Errorf("Hashtable didn't get the correct size.")
	}
	fmt.Printf("TestSize finished running with status OK\n")
}
