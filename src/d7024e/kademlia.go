package d7024e

type Kademlia struct {
	network      *Network
	routingTable *RoutingTable
	alpha        int // also known as the alpha value that determines how many concurrent findclosestcontacts calls will exist
}

func (kademlia *Kademlia) LookupContact(target *Contact) Contact {
	c1 := make(chan []Contact)
	//	c2 := make(chan []Contact)
	var initiatorList []Contact
	var candidateList ContactCandidates
	//var ContactedList []Contact
	initiatorList = kademlia.routingTable.FindClosestContacts(target.ID, kademlia.alpha)
	for i := 0; i < kademlia.alpha; i++ {
		go kademlia.GoFindNode(target, &initiatorList[i], c1)
	}
	candidateList.Append(<-c1)
	candidateList.Sort()
	return candidateList.GetContacts(1)[0]
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) channelRec(contact *Contact, channel chan []Contact) {

	//hannel <- kademlia.network.SendFindContactMessage(contact)
}

func (kademlia *Kademlia) GoFindNode(target *Contact, contact *Contact, channel chan []Contact) {
	var queriedList []Contact
	var requestList ContactCandidates
	var resultList = kademlia.network.SendFindContactMessage(target, contact)
	var flag = true
	for ok := true; ok; ok = flag {
		for i := 0; i < len(resultList); i++ {
			if kademlia.EqualKademliaID(queriedList, &resultList[i]) {
				queriedList[i] = resultList[i]
				requestList.Append(kademlia.network.SendFindContactMessage(target, &resultList[i]))
			}
		}
		requestList.Sort()
		var tempCon Contact
		resultList[0].CalcDistance(target.ID)
		tempCon = requestList.GetContacts(1)[0]
		tempCon.CalcDistance(target.ID)
		if !(tempCon.Less(&resultList[0])) {
			ok = false
		} else {
			resultList = requestList.GetContacts(bucketSize)
		}
	}
	channel <- resultList
}

// Searches argument list to see whether or not argument contact exists
func (kademlia *Kademlia) EqualKademliaID(contactList []Contact, contact *Contact) bool {
	for i := 0; i < len(contactList); i++ {
		if contact.ID == contactList[i].ID {
			return true
		}
	}
	return false
}
