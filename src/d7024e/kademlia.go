package d7024e

import (
	"encoding/hex"
	"fmt"
	"time"
)

// Kademlia type definition
type Kademlia struct {
	rt    			*RoutingTable
	ht				*ValueHashtable
	alpha 			int
	k 	  			int
	maxRoundTime 	time.Duration
}

const k = 20
const alpha = 3
const maxRoundTime = 3 * time.Second

// NewKademlia Constructor function for Kademlia class
func NewKademlia(rt *RoutingTable, ht *ValueHashtable) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.rt = rt
	kademlia.ht = ht
	kademlia.alpha = alpha
	kademlia.k = k
	kademlia.maxRoundTime = maxRoundTime
	return kademlia
}

// LookupContact returns the closest found contact of the searched network based on given target
func (kademlia *Kademlia) LookupContact(target *Contact, count int, findNodeRequestChannel chan findNodeRequest) []Contact {
	start := time.Now()
	countClosestIndex := count - 1
	var countClosestContact *Contact
	// Add the node who performs the lookup to the shortlist and set it to queried so it does not query itself
	shortlist := NewContactCandidates([]Contact{*kademlia.rt.me})
	shortlist.contacts[0].CalcDistance(target.ID)
	shortlist.contacts[0].queried = true
	// Find alpha closest contacts of own buckets
	shortlist.Append(kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha))
	// Remove any potential duplicates and sort by distance to target (if FindClosestContacts returned the node who is performing the lookup there would be a duplicate)
	shortlist.RemoveDuplicates()
	shortlist.Sort()
	// Note closest contact
	countClosestContact = &shortlist.contacts[MinInt(countClosestIndex, shortlist.Len() - 1)]
	// If the distance to the closest contact is 0 the lookup is done (Exit condition)
	if !countClosestContact.Distance.EqualsZero() {
		numQueriedContacts := 0
		queryNewContacts := true
		// If we have queried k contacts the lookup is done (Exit condition)
		for numQueriedContacts < kademlia.k {
			// Find contacts to query that has not already been queried
			currentQueryContacts := kademlia.findQueryContacts(shortlist)
			// If there are no contacts to be queried the lookup is done (Exit condition)
			if len(currentQueryContacts) <= 0 {
				break
			}
			findNodeResponseChannel := make(chan findNodeResponse, len(currentQueryContacts))
			roundStartTime := time.Now()
			kademlia.sendFindNodeRequests(target, currentQueryContacts, findNodeRequestChannel, findNodeResponseChannel)
			numQueriedContacts += kademlia.receiveFindNodeResponses(&currentQueryContacts, shortlist, roundStartTime, queryNewContacts, findNodeResponseChannel)
			// Remove contacts from shortlist that hasn't responded yet
			shortlist.removeContacts(currentQueryContacts)
			// If a duplicate contact of an already queried contact has been added, the duplicate will most likely
			// have the queried variable set to false. This should not be a problem since we remove every duplicate
			// of every contact expect the first which should be the one with the queried variable set to true if it
			// has been queried already
			shortlist.RemoveDuplicates()
			shortlist.Sort()
			// We only need to check if the newCountClosestContact is closer than countClosestContact if we have
			// a total of count contacts in our shortlist
			if shortlist.Len() >= count {
				newCountClosestContact := &shortlist.contacts[MinInt(countClosestIndex, shortlist.Len() - 1)]
				// If a cycle doesn't find a closer contact we don't want new found contacts to be queried
				if !newCountClosestContact.Distance.Less(countClosestContact.Distance) {
					queryNewContacts = false
				}
			} else {
				countClosestContact = &shortlist.contacts[MinInt(countClosestIndex, shortlist.Len() - 1)]
			}
			// If the distance to the closest contact is 0 the lookup is done (Exit condition)
			if countClosestContact.Distance.EqualsZero() {
				break
			}
		}
	}
	fmt.Println("Lookup contact took", time.Since(start))
	return shortlist.contacts[:MinInt(count, shortlist.Len())]
}

// JoinNetwork attempts to join a network given the address of a participant in an existing network
func (kademlia *Kademlia) JoinNetwork(address string, findNodeRequestChannel chan findNodeRequest) {
	// First send a FIND_NODE RPC call to the bootstrap node with self as target, the ID of the bootstrap node should not matter since it is not used when performing the lookup
	findNodeResponseChannel := make(chan findNodeResponse)
	findNodeRequestChannel <- findNodeRequest{kademlia.rt.me, NewContact(NewRandomKademliaID(), address), findNodeResponseChannel}
	findNodeResponse := <- findNodeResponseChannel
	// Add each of the returned contacts from the bootstrap node to the joining node's buckets
	for i := 0; i < len(findNodeResponse.contacts); i++ {
		kademlia.rt.AddContact(findNodeResponse.contacts[i])
	}
	// Perform an refresh by executing a lookup with a random random ID as target which is not the same as the joining node's or it's neighbour
	refreshContact := NewContact(kademlia.rt.me.ID.NewKademliaIDWithinRange(), "") // Probably not necessary to check for the collisions of the new random ID due to the large ID "space"
	kademlia.LookupContact(refreshContact, kademlia.k, findNodeRequestChannel)
}

func (kademlia *Kademlia) LookupData(hash [HashSize]byte, findDataRequestChannel chan findDataRequest) []byte {
	start := time.Now()
	data := kademlia.ht.Get(hash)
	// Searching node already has hash in hashtable (Exit condition)
	if data == nil {
		kademliaID, _ := NewKademliaID(hex.EncodeToString(hash[:]))
		target := NewContact(kademliaID, "")
		// Add the node who performs the lookup to the shortlist and set it to queried so it does not query itself
		shortlist := NewContactCandidates([]Contact{*kademlia.rt.me})
		shortlist.contacts[0].CalcDistance(target.ID)
		shortlist.contacts[0].queried = true
		// Find alpha closest contacts of own buckets
		shortlist.Append(kademlia.rt.FindClosestContacts(target.ID, kademlia.alpha))
		// Remove any potential duplicates and sort by distance to target (if FindClosestContacts returned the node who is performing the lookup there would be a duplicate)
		shortlist.RemoveDuplicates()
		shortlist.Sort()
		// Note closest contact
		closestContact := &shortlist.contacts[0]
		numQueriedContacts := 0
		queryNewContacts := true
		// If we have queried k contacts the lookup is done (Exit condition)
		for numQueriedContacts < kademlia.k {
			// Find contacts to query that has not already been queried
			currentQueryContacts := kademlia.findQueryContacts(shortlist)
			// If there are no contacts to be queried the lookup is done (Exit condition)
			if len(currentQueryContacts) <= 0 {
				break
			}
			findDataResponseChannel := make(chan findDataResponse, len(currentQueryContacts))
			roundStartTime := time.Now()
			kademlia.sendFindDataRequests(hash, target, currentQueryContacts, findDataRequestChannel, findDataResponseChannel)
			numResponses := 0
			data, numResponses = kademlia.receiveFindDataResponses(&currentQueryContacts, shortlist, roundStartTime, queryNewContacts, findDataResponseChannel)
			if data != nil {
				// Data was found (Exit condition)
				break
			}
			numQueriedContacts += numResponses
			// Remove contacts from shortlist that hasn't responded yet
			shortlist.removeContacts(currentQueryContacts)
			// If a duplicate contact of an already queried contact has been added, the duplicate will most likely
			// have the queried variable set to false. This should not be a problem since we remove every duplicate
			// of every contact expect the first which should be the one with the queried variable set to true if it
			// has been queried already
			shortlist.RemoveDuplicates()
			shortlist.Sort()
			newCountClosestContact := &shortlist.contacts[0]
			// If a cycle doesn't find a closer contact we don't want new found contacts to be queried
			if !newCountClosestContact.Distance.Less(closestContact.Distance) {
				queryNewContacts = false
			}
		}
	}
	fmt.Println("Lookup data took", time.Since(start))
	return data
}
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
*/
/*
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
func (kademlia *Kademlia) findQueryContacts(shortlist *ContactCandidates) []*Contact{
	var queryContacts []*Contact;
	for i := 0; i < shortlist.Len(); i++ {
		if !shortlist.contacts[i].queried {
			queryContacts = append(queryContacts, &shortlist.contacts[i])
			if len(queryContacts) >= kademlia.alpha {
				break
			}
		}
	}
	return queryContacts
}

func (kademlia *Kademlia) sendFindNodeRequests(target *Contact, queryContacts []*Contact, findNodeRequestChannel chan findNodeRequest, findNodeResponseChannel chan findNodeResponse) {
	for i := 0; i < MinInt(alpha, len(queryContacts)); i++ {
		findNodeRequestChannel <- findNodeRequest{target, queryContacts[i], findNodeResponseChannel}
	}
}

func (kademlia *Kademlia) receiveFindNodeResponses(queryContacts *[]*Contact, shortlist *ContactCandidates, roundStartTime time.Time, queryNewContacts bool, findNodeResponseChannel chan findNodeResponse) int {
	responsesReceived := 0
	expectedResponses := len(*queryContacts)
	// Wait either for kademlia.maxRoundTime or until we have gotten all expected responses depending on
	// which is the quickest and handle every response we get
	for time.Since(roundStartTime) < kademlia.maxRoundTime && responsesReceived < expectedResponses {
		select {
		case findNodeResponse, ok := <- findNodeResponseChannel:
			if ok {
				// TODO handle replies later than maxRoundTime
				findNodeResponse.sender.queried = true
				if !queryNewContacts {
					// Since we don't want to query any of the new contacts we set their queried variable to true
					for i := 0; i < len(findNodeResponse.contacts); i++ {
						findNodeResponse.contacts[i].queried = true
					}
				}
				shortlist.Append(findNodeResponse.contacts)
				responsesReceived++
				// Since we got a response we remove the responder from the currentQueryContacts
				remaininqQueryContacts := removeContact(*queryContacts, findNodeResponse.sender.ID)
				(*queryContacts) = remaininqQueryContacts
			} else {
				fmt.Println("\"findNodeResponseChannel\" has been closed!")
			}
		default:
		}
	}
	return responsesReceived
}

func (kademlia *Kademlia) sendFindDataRequests(hash [HashSize]byte, target *Contact, queryContacts []*Contact, findDataRequestChannel chan findDataRequest, findDataResponseChannel chan findDataResponse) {
	for i := 0; i < MinInt(alpha, len(queryContacts)); i++ {
		findDataRequestChannel <- findDataRequest{hash, target, queryContacts[i], findDataResponseChannel}
	}
}

func (kademlia *Kademlia) receiveFindDataResponses(queryContacts *[]*Contact, shortlist *ContactCandidates, roundStartTime time.Time, queryNewContacts bool, findDataResponseChannel chan findDataResponse) ([]byte, int) {
	responsesReceived := 0
	expectedResponses := len(*queryContacts)
	// Wait either for kademlia.maxRoundTime or until we have gotten all expected responses depending on
	// which is the quickest and handle every response we get
	for time.Since(roundStartTime) < kademlia.maxRoundTime && responsesReceived < expectedResponses {
		select {
		case findDataResponse, ok := <- findDataResponseChannel:
			if ok {
				responsesReceived++
				if findDataResponse.data != nil {
					// Data was found (Exit condition)
					return findDataResponse.data, responsesReceived
				}
				// TODO handle replies later than maxRoundTime
				findDataResponse.sender.queried = true
				if !queryNewContacts {
					// Since we don't want to query any of the new contacts we set their queried variable to true
					for i := 0; i < len(findDataResponse.contacts); i++ {
						findDataResponse.contacts[i].queried = true
					}
				}
				shortlist.Append(findDataResponse.contacts)
				// Since we got a response we remove the responder from the currentQueryContacts
				remaininqQueryContacts := removeContact(*queryContacts, findDataResponse.sender.ID)
				(*queryContacts) = remaininqQueryContacts
			} else {
				fmt.Println("\"findDataResponseChannel\" has been closed!")
			}
		default:
		}
	}
	return nil, responsesReceived
}