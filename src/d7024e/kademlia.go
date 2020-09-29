package d7024e

import (
	"fmt"
	"time"
	"encoding/hex"
)

type Kademlia struct {
	net   *Network
	rt    *RoutingTable
	alpha int
}

// NewKademlia Constructor function for Kademlia class
func NewKademlia(net *Network, rt *RoutingTable, alpha int) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.net = net
	kademlia.rt = rt
	kademlia.alpha = alpha
	return kademlia
}

func printContacts(contacts []Contact) {
	for i, contact := range contacts {
		fmt.Println(i, contact.String(), "distance:", contact.Distance)
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) Contact {
	start := time.Now()
	c1 := make(chan []Contact)
	var initiatorList []Contact
	var candidateList ContactCandidates
	initiatorList = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	var closestContact Contact
	if initiatorList[0].Distance.EqualsZero() {
		closestContact = initiatorList[0]
	} else {
		for i := 0; i < MinInt(kademlia.alpha, len(initiatorList)); i++ {
			go kademlia.goFindNode(target, &initiatorList[i], c1)
		}
		for i := 0; i < MinInt(kademlia.alpha, len(initiatorList)); i++ {
			newCandidates := <-c1
			candidateList.Append(newCandidates)
			if newCandidates[0].Distance.EqualsZero() {
				break
			}
		}
		candidateList.Sort()
		closestContact = candidateList.GetContacts(1)[0]
	}
	fmt.Println("Lookup contact took", time.Since(start))
	return closestContact
}

func (kademlia *Kademlia) LookupData(hash [HashSize]byte) []byte {
	start := time.Now()
	data := kademlia.net.ht.Get(hash)
	if (data == nil) {
		c1 := make(chan []byte)
		var initiatorList []Contact
		target := NewContact(NewKademliaID(hex.EncodeToString(hash[:])), "")
		initiatorList = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
		for i := 0; i < MinInt(kademlia.alpha, len(initiatorList)); i++ {
			go kademlia.goFindData(hash, &initiatorList[i], c1)
		}
		data = <-c1
	}
	fmt.Println("Lookup data took", time.Since(start))
	return data
}

func (kademlia *Kademlia) Store(filename []byte, data []byte) [HashSize]byte{
	start := time.Now()
	hash := Hash(filename)
	target := NewContact(NewKademliaID(hex.EncodeToString(hash[:])), "")
	var initiatorList []Contact
	initiatorList = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	var candidateList ContactCandidates
	c1 := make(chan []Contact)
	for i := 0; i < MinInt(kademlia.alpha, len(initiatorList)); i++ {
		go kademlia.goFindNode(target, &initiatorList[i], c1)
	}
	for i := 0; i < MinInt(kademlia.alpha, len(initiatorList)); i++ {
		newCandidates := <-c1
		candidateList.Append(newCandidates)
	}
	candidateList.RemoveDuplicates()
	candidateList.Sort()
	closestContacts := candidateList.GetContacts(MinInt(IDLength, candidateList.Len()))
	for _, contact := range closestContacts {
		kademlia.net.SendStoreMessage(&contact, hash, data)
	}
	fmt.Println("Store data took", time.Since(start))
	return hash
}

func (kademlia *Kademlia) goFindNode(target *Contact, contact *Contact, channel chan []Contact) {
	queriedList := []Contact{*kademlia.rt.me}
	var requestList ContactCandidates
	var resultList = kademlia.net.SendFindContactMessage(target, contact)
	if !resultList[0].Distance.EqualsZero() {
		var flag = true
		for ok := true; ok; ok = flag {
			for i := 0; i < len(resultList); i++ {
				if !kademlia.EqualKademliaID(queriedList, &resultList[i]) {
					queriedList = append(queriedList, resultList[i])
					currentResponseList := kademlia.net.SendFindContactMessage(target, &resultList[i])
					requestList.Append(currentResponseList)
					if currentResponseList[0].Distance.EqualsZero() {
						flag = false
						break
					}
				}
			}
			if requestList.contacts != nil {
				worstResult := resultList[len(resultList) - 1]
				mergeList := requestList
				mergeList.Append(resultList)
				mergeList.RemoveDuplicates()
				mergeList.Sort()
				worstMergeMaxAllowed := MinInt(IDLength, mergeList.Len())
				if !mergeList.GetContacts(worstMergeMaxAllowed)[worstMergeMaxAllowed - 1].Distance.Less(worstResult.Distance) && len(resultList) >= IDLength {
					flag = false
				} else {
					resultList = mergeList.GetContacts(worstMergeMaxAllowed)
				}
			} else {
				flag = false
			}
		}
	}
	channel <- resultList
}

func (kademlia *Kademlia) goFindData(hash [HashSize]byte, contact *Contact, channel chan []byte) {
	data := kademlia.net.SendFindDataMessage(hash, contact)
	if (data == nil) {
		target := NewContact(NewKademliaID(hex.EncodeToString(hash[:])), "")
		queriedList := []Contact{*kademlia.rt.me}
		var requestList ContactCandidates
		var resultList = kademlia.net.SendFindContactMessage(target, contact)
		var flag = true
		for ok := true; ok; ok = flag {
			for i := 0; i < len(resultList); i++ {
				if !kademlia.EqualKademliaID(queriedList, &resultList[i]) {
					queriedList = append(queriedList, resultList[i])
					data = kademlia.net.SendFindDataMessage(hash, contact)
					if (data == nil) {
						currentResponseList := kademlia.net.SendFindContactMessage(target, &resultList[i])
						requestList.Append(currentResponseList)
					} else {
						flag = false
						break
					}
				}
			}
			if requestList.contacts != nil && flag {
				worstResult := resultList[len(resultList) - 1]
				mergeList := requestList
				mergeList.Append(resultList)
				mergeList.RemoveDuplicates()
				mergeList.Sort()
				worstMergeMaxAllowed := MinInt(IDLength, mergeList.Len())
				if !mergeList.GetContacts(worstMergeMaxAllowed)[worstMergeMaxAllowed - 1].Distance.Less(worstResult.Distance) && len(resultList) >= IDLength {
					flag = false
				} else {
					resultList = mergeList.GetContacts(worstMergeMaxAllowed)
				}
			} else {
				flag = false
			}
		}
	}
	channel <- data
}

// Searches argument list to see whether or not argument contact exists
func (kademlia *Kademlia) EqualKademliaID(contactList []Contact, contact *Contact) bool {
	for i := 0; i < len(contactList); i++ {
		if contact.ID.Equals(contactList[i].ID) {
			return true
		}
	}
	return false
}

func MinInt(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}
