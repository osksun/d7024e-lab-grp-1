package d7024e

import (
	"fmt"
)

type Kademlia struct {
    net		*Network
    rt		*RoutingTable
    alpha	int // also known as the alpha value that determines how many concurrent findclosestcontacts calls will exist
}


func printContacts(contacts []Contact) {
	for i, contact := range contacts {
		fmt.Println(i, contact.String(), "distance:", contact.distance)
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	//c1 := make(chan []Contact)
	var initiatorList []Contact
	//var uncontactedContactsList []Contact
	initiatorList = kademlia.rt.FindClosestContacts(target.ID, bucketSize)
	printContacts(initiatorList)
	/*
	for i := 0; i < kademlia.alpha; i++ {
		go kademlia.channelRec(&initiatorList[i], c1)
	}
	*/
	/*
	var c1Output []Contact = <-c1
	for i := 0; i < kademlia.alpha; i++ {
		go kademlia.channelRec(&c1Output[i], c1)
	}
	*/
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) channelRec(contact *Contact, channel chan []Contact) {

	//channel <- kademlia.network.SendFindContactMessage(contact)
}