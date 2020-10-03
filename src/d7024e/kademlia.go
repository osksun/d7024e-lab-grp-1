package d7024e

import (
	"encoding/hex"
	"fmt"
	"time"
)

type Kademlia struct {
	rt    			*RoutingTable
	alpha 			int
	k 	  			int
	maxRoundTime 	time.Duration
}

const k = 20
const alpha = 3
const maxRoundTime = 10 * time.Second

// NewKademlia Constructor function for Kademlia class
func NewKademlia(rt *RoutingTable) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.rt = rt
	kademlia.alpha = alpha
	kademlia.k = k
	kademlia.maxRoundTime = maxRoundTime
	return kademlia
}

// LookupContact returns the closest found contact of the searched network based on given target
func (kademlia *Kademlia) LookupContact(target *Contact, findNodeRequestChannel chan findNodeRequest) Contact {
	start := time.Now()
	var closestContact Contact
	// Add the node who performs the lookup to the shortlist and set it to queried so it does not query itself
	shortlist := NewContactCandidates([]Contact{*kademlia.rt.me})
	shortlist.contacts[0].CalcDistance(target.ID)
	shortlist.contacts[0].queried = true
	// Find alpha closest contacts of own buckets
	shortlist.Append(kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha))
	// Remove any potential duplicates (if FindClosestContacts returned the node who is performing the lookup there would be a duplicate)
	shortlist.RemoveDuplicates()
	shortlist.Sort()
	// Note closest contact
	closestContact = shortlist.contacts[0]
	numQueriedContacts := 0
	queryNewContacts := true
	// If the distance to the closest contact is 0 the lookup is done (Exit condition)
	if !closestContact.Distance.EqualsZero() {
		// If we have queried k contacts the lookup is done (Exit condition)
		for numQueriedContacts < kademlia.k {
			// Find contacts to query that has not already been queried
			var currentQueryContacts []*Contact
			for i := 0; i < shortlist.Len(); i++ {
				if !shortlist.contacts[i].queried {
					currentQueryContacts = append(currentQueryContacts, &shortlist.contacts[i])
					if len(currentQueryContacts) >= alpha {
						break
					}
				}
			}
			// If there are no contacts to be queried the lookup is done (Exit condition)
			if len(currentQueryContacts) <= 0 {
				break
			}
			findNodeResponseChannel := make(chan findNodeResponse, len(currentQueryContacts))
			roundStartTime := time.Now()
			// Send parallel, asynchronous FIND_NODE RPCs to at most alpha contacts in the shortlist
			for i := 0; i < MinInt(alpha, len(currentQueryContacts)); i++ {
				findNodeRequestChannel <- findNodeRequest{target, currentQueryContacts[i], findNodeResponseChannel}
			}
			// Wait either for kademlia.maxRoundTime or untill we have gotten all expected responces depending on
			// which is the quickest and handle every response we get
			nResponses := 0
			nExpectedResponses := len(currentQueryContacts)
			for time.Since(roundStartTime) < kademlia.maxRoundTime && nResponses < nExpectedResponses {
				select {
				case RPCResponse, ok := <- findNodeResponseChannel:
					if ok {
						// TODO handle replies later than maxRoundTime
						RPCResponse.sender.queried = true
						if !queryNewContacts {
							// Since we don't want to query any of the new contacts we set their queried variable to true
							for i := 0; i < len(RPCResponse.contacts); i++ {
								RPCResponse.contacts[i].queried = true
							}
						}
						shortlist.Append(RPCResponse.contacts)
						nResponses++
						// Since we got a response we remove the responder from the currentQueryContacts
						currentQueryContacts = removeContact(currentQueryContacts, RPCResponse.sender.ID)
					} else {
						fmt.Println("\"findNodeResponseChannel\" has been closed!")
					}
				default:
				}
			}
			// Remove contacts from shortlist that hasn't responded yet
			for _, contact := range currentQueryContacts {
				shortlist.remove(contact.ID)
			}
			// Update the number of queried contacts
			numQueriedContacts += nResponses
			// If a duplicate contact of an already queried contact has been added, the duplicate will most likely
			// have the queried variable set to false. This should not be a problem since we remove every duplicate
			// of every contact expect the first which should be the one with the queried variable set to true if it
			// has been queried already
			shortlist.RemoveDuplicates()
			shortlist.Sort()
			// If a cycle doesn't find a closer contact we don't want new found contacts to be queried
			if !shortlist.contacts[0].Distance.Less(closestContact.Distance) {
				queryNewContacts = false
			}
			// Note closest contact
			closestContact = shortlist.contacts[0]
			// If the distance to the closest contact is 0 the lookup is done (Exit condition)
			if closestContact.Distance.EqualsZero() {
				break
			}
		}
	}
	fmt.Println("Lookup contact took", time.Since(start))
	return closestContact
}

func (kademlia *Kademlia) LookupData(hash [HashSize]byte) []byte {
	start := time.Now()
	data := kademlia.net.ht.Get(hash)
	if data == nil {
		c1 := make(chan []byte)
		var initiatorList []Contact
		kID, _ := NewKademliaID(hex.EncodeToString(hash[:]))
		target := NewContact(kID, "")
		initiatorList = kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha)
		for i := 0; i < MinInt(kademlia.alpha, len(initiatorList)); i++ {
			go kademlia.goFindData(hash, &initiatorList[i], c1)
		}
		data = <-c1
	}
	fmt.Println("Lookup data took", time.Since(start))
	return data
}

func (kademlia *Kademlia) Store(filename []byte, data []byte) [HashSize]byte {
	start := time.Now()
	hash := Hash(filename)
	kID, _ := NewKademliaID(hex.EncodeToString(hash[:]))
	target := NewContact(kID, "")
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
	var resultList = kademlia.net.SendFindContactMessage(target, contact)
	if !resultList[0].Distance.EqualsZero() {
		var flag = true
		for ok := true; ok; ok = flag {
			var requestList ContactCandidates
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
				worstResult := resultList[len(resultList)-1]
				mergeList := requestList
				mergeList.Append(resultList)
				mergeList.RemoveDuplicates()
				mergeList.Sort()
				worstMergeMaxAllowed := MinInt(IDLength, mergeList.Len())
				if !mergeList.GetContacts(worstMergeMaxAllowed)[worstMergeMaxAllowed-1].Distance.Less(worstResult.Distance) && len(resultList) >= IDLength {
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
	if data == nil {
		kID, _ := NewKademliaID(hex.EncodeToString(hash[:]))
		target := NewContact(kID, "")
		queriedList := []Contact{*kademlia.rt.me}
		var requestList ContactCandidates
		var resultList = kademlia.net.SendFindContactMessage(target, contact)
		var flag = true
		for ok := true; ok; ok = flag {
			for i := 0; i < len(resultList); i++ {
				if !kademlia.EqualKademliaID(queriedList, &resultList[i]) {
					queriedList = append(queriedList, resultList[i])
					data = kademlia.net.SendFindDataMessage(hash, contact)
					if data == nil {
						currentResponseList := kademlia.net.SendFindContactMessage(target, &resultList[i])
						requestList.Append(currentResponseList)
					} else {
						flag = false
						break
					}
				}
			}
			if requestList.contacts != nil && flag {
				worstResult := resultList[len(resultList)-1]
				mergeList := requestList
				mergeList.Append(resultList)
				mergeList.RemoveDuplicates()
				mergeList.Sort()
				worstMergeMaxAllowed := MinInt(IDLength, mergeList.Len())
				if !mergeList.GetContacts(worstMergeMaxAllowed)[worstMergeMaxAllowed-1].Distance.Less(worstResult.Distance) && len(resultList) >= IDLength {
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
