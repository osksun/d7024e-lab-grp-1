package d7024e

import (
	//"container/list"
	"fmt"
	"testing"
)

// bucket definition
// contains a List
// type bucket struct {
// 	list *list.List
// }

func TestNewBucket(t *testing.T) {
	bucket := newBucket()
	lType := fmt.Sprintf("%T", bucket.list)
	bType := fmt.Sprintf("%T", bucket)
	if lType != "*list.List" {
		t.Errorf("The bucket list is not of type list")
	}
	if bType != "*d7024e.bucket" {
		t.Errorf("The bucket is not of type bucket")
	}
	fmt.Printf("TestNewBucket finished running with status OK\n")
}

func TestAddContact(t *testing.T) {
	nct := *NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	nct2 := *NewContact(NewKademliaID("FFFFFFF000000000000000000000000000000000"), "localhost:8002")
	bucket := newBucket()
	ls1 := bucket.Len()
	bucket.AddContact(nct)
	ls2 := bucket.Len()
	if !(ls1 < ls2) {
		t.Errorf("Bucket size didn't increase when adding as it should.")
	}
	bucket.AddContact(nct2)
	bucket.AddContact(nct)
	if bucket.GetFirst().ID != nct.ID {
		t.Errorf("Bucket didn't move the contact to the front.")
	}
	fmt.Printf("TestAddContact finished running with status OK\n")
}

func TestGetContactAndCalcDistance(t *testing.T) {
	nct := *NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	nct2 := *NewContact(NewKademliaID("FFFFFFF000000000000000000000000000000000"), "localhost:8002")
	bucket := newBucket()
	bucket.AddContact(nct)
	bucket.AddContact(nct2)
	kdid := NewKademliaID("0000000000000000000000000000000000000000")
	contacts := bucket.GetContactAndCalcDistance(kdid)

	if len(contacts) != 2 {
		t.Errorf("Bucket didn't return all of the contacts")
	}
	if contacts[0].Distance.String() != "fffffff000000000000000000000000000000000" || contacts[1].Distance.String() != "ffffffff00000000000000000000000000000000" {
		t.Errorf("Bucket didn't calculate the distance correctly.")
	}
	fmt.Printf("TestGetContactAndCalcDistance finished running with status OK\n")
}

// // GetContactAndCalcDistance returns an array of Contacts where
// // the distance has already been calculated
// func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
// 	var contacts []Contact

// 	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
// 		contact := elt.Value.(Contact)
// 		contact.CalcDistance(target)
// 		contacts = append(contacts, contact)
// 	}

// 	return contacts
// }

// // Len return the size of the bucket
// func (bucket *bucket) Len() int {
// 	return bucket.list.Len()
// }

// func (bucket *bucket) GetLast() Contact {
// 	var contact Contact
// 	contact = bucket.list.Back().Value.(Contact)
// 	return contact
// }
