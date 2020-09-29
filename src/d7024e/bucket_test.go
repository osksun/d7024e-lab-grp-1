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
}

// }
// // newBucket returns a new instance of a bucket
// func newBucket() *bucket {
// 	bucket := &bucket{}
// 	bucket.list = list.New()
// 	return bucket
// }

// // AddContact adds the Contact to the front of the bucket
// // or moves it to the front of the bucket if it already existed
// func (bucket *bucket) AddContact(contact Contact) {
// 	var element *list.Element
// 	for e := bucket.list.Front(); e != nil; e = e.Next() {
// 		nodeID := e.Value.(Contact).ID

// 		if (contact).ID.Equals(nodeID) {
// 			element = e
// 		}
// 	}

// 	if element == nil {
// 		if bucket.list.Len() < bucketSize {
// 			bucket.list.PushFront(contact)
// 		}
// 	} else {
// 		bucket.list.MoveToFront(element)
// 	}
// }

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
