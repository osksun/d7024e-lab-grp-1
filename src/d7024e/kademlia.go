package d7024e

import (
	"fmt"
	"time"
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
		fmt.Println(i, contact.String(), "distance:", contact.distance)
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) Contact {
	start := time.Now()
	c1 := make(chan []Contact)
	//	c2 := make(chan []Contact)
	var initiatorList []Contact
	var candidateList ContactCandidates
	//var ContactedList []Contact
	initiatorList = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
	for i := 0; i < kademlia.alpha; i++ {
		go kademlia.goFindNode(target, &initiatorList[i], c1)
	}
	for i := 0; i < kademlia.alpha; i++ {
		candidateList.Append(<-c1)
	}
	candidateList.Sort()
	fmt.Println("Lookup took", time.Since(start))
	return candidateList.GetContacts(1)[0]
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) goFindNode(target *Contact, contact *Contact, channel chan []Contact) {
	var queriedList []Contact
	var requestList ContactCandidates
	var resultList = kademlia.net.SendFindContactMessage(target, contact)
	var flag = true
	for ok := true; ok; ok = flag {
		for i := 0; i < len(resultList); i++ {
			if !kademlia.EqualKademliaID(queriedList, &resultList[i]) {
				queriedList = append(queriedList, resultList[i])
				requestList.Append(kademlia.net.SendFindContactMessage(target, &resultList[i]))
			}
		}
		for i := 0; i < requestList.Len(); i++ {
			requestList.contacts[i].CalcDistance(target.ID)
		}
		requestList.Sort()
		var tempCon Contact
		for i := 0; i < len(resultList); i++ {
			resultList[i].CalcDistance(target.ID)
		}
		if requestList.contacts != nil {
			tempCon = requestList.GetContacts(1)[0]
			tempCon.CalcDistance(target.ID)
			if !(tempCon.Less(&resultList[0])) {
				flag = false
			} else {
				resultList = requestList.GetContacts(requestList.Len())
			}
		} else {
			flag = false
		}
	}
	channel <- resultList
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
