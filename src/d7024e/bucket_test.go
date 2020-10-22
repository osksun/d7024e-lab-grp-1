package d7024e

import (
	"fmt"
	"testing"
)

func TestNewBucket(t *testing.T) {
	bucket := newBucket()
	lType := fmt.Sprintf("%T", bucket.list)
	bType := fmt.Sprintf("%T", bucket)
	if lType != "*list.List" {
		t.Error("The bucket list is not of type list")
	}
	if bType != "*d7024e.bucket" {
		t.Error("The bucket is not of type bucket")
	}
	//fmt.Println("TestNewBucket finished running with status OK")
}

func TestBucketAddContact(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct := *NewContact(kID1, "localhost:8001")
	nct2 := *NewContact(kID2, "localhost:8002")
	bucket := newBucket()
	ls1 := bucket.Len()
	bucket.AddContact(nct)
	ls2 := bucket.Len()
	if !(ls1 < ls2) {
		t.Error("Bucket size didn't increase when adding as it should.")
	}
	bucket.AddContact(nct2)
	bucket.AddContact(nct)
	if bucket.GetFirst().ID != nct.ID {
		t.Error("Bucket didn't move the contact to the front.")
	}
	//fmt.Println("TestBucketAddContact finished running with status OK")
}

func TestGetContactAndCalcDistance(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct := *NewContact(kID1, "localhost:8001")
	nct2 := *NewContact(kID2, "localhost:8002")
	bucket := newBucket()
	bucket.AddContact(nct)
	bucket.AddContact(nct2)
	kdid, _ := NewKademliaID("0000000000000000000000000000000000000000")
	contacts := bucket.GetContactAndCalcDistance(kdid)

	if len(contacts) != 2 {
		t.Error("Bucket didn't return all of the contacts")
	}
	if contacts[0].Distance.String() != "fffffff000000000000000000000000000000000" || contacts[1].Distance.String() != "ffffffff00000000000000000000000000000000" {
		t.Error("Bucket didn't calculate the distance correctly.")
	}
	//fmt.Println("TestGetContactAndCalcDistance finished running with status OK")
}

func TestLen(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct := *NewContact(kID1, "localhost:8001")
	nct2 := *NewContact(kID2, "localhost:8002")
	bucket := newBucket()
	bucket.AddContact(nct)
	bucket.AddContact(nct2)
	if bucket.Len() != 2 {
		t.Error("Bucket didn't return the correct length.")
	}
	//fmt.Println("TestLen finished running with status OK")
}

func TestGetLast(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct := *NewContact(kID1, "localhost:8001")
	nct2 := *NewContact(kID2, "localhost:8002")
	bucket := newBucket()
	bucket.AddContact(nct)
	bucket.AddContact(nct2)

	if bucket.GetLast().ID != nct.ID {
		t.Error("Bucket didn't get the last correct element.")
	}
	//fmt.Println("TestGetLast finished running with status OK")
}

func TestGetFirst(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct := *NewContact(kID1, "localhost:8001")
	nct2 := *NewContact(kID2, "localhost:8002")
	bucket := newBucket()
	bucket.AddContact(nct)
	bucket.AddContact(nct2)

	if bucket.GetFirst().ID != nct2.ID {
		t.Error("Bucket didn't get the first correct element.")
	}
	//fmt.Println("TestGetFirst finished running with status OK")
}
