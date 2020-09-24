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
		fmt.Println(i, contact.String(), "distance:", contact.Distance)
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
	fmt.Println("Lookup took", time.Since(start))
	return closestContact
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
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
					if currentResponseList[0].Distance.EqualsZero() {
						// write to channel return
					}
					requestList.Append(currentResponseList)
				}
			}
			requestList.Sort()

			if requestList.contacts != nil {
				closestCandidate := requestList.GetContacts(1)[0]
				if !(closestCandidate.Less(&resultList[0])) {
					flag = false
				} else {
					resultList = requestList.GetContacts(requestList.Len())
				}
			} else {
				flag = false
			}
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

func MinInt(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}
