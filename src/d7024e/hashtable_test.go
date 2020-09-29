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
		t.Error("The hashtable is not of type ValueHashtable")
	}
	fmt.Println("TestNewValueHashtable finished running with status OK")
}

func TestHash(t *testing.T) {
	k1 := []byte("key")
	k2 := []byte("key")
	h1 := Hash(k1)
	h2 := Hash(k2)
	if h1 != h2 {
		t.Error("Hash produced is not consistent with input.")
	}
	fmt.Println("TestHash finished running with status OK")
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
		t.Error("Hashtable items wasn't initialized as nil or not made.")
	}

	testitems := make(map[[HashSize]byte][]byte)
	testitems[k] = v1

	if !reflect.DeepEqual(ht.items, testitems) {
		t.Error("Hashtable didn't put correctly.")
	}
	fmt.Println("TestPut finished running with status OK")
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
		t.Error("Hashtable didn't remove correctly.")
	}
	fmt.Println("TestRemove finished running with status OK")
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
		t.Error("Hashtable didn't get correctly.")
	}
	fmt.Println("TestGet finished running with status OK")
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
		t.Error("Hashtable didn't get the correct size.")
	}
	fmt.Println("TestSize finished running with status OK")
}
