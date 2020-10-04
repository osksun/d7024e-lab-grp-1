package d7024e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Network struct {
	rt				*RoutingTable
	ht				*ValueHashtable
	findNodeChannel chan findNodeRequest
	exitChannel		chan bool
}

type msg struct {
	Message string
	Hash    [HashSize]byte
	Data    []byte
	Target  Contact
	Sender  Contact
}

type response_msg struct {
	Message     string
	ContactList []Contact
	Data        []byte
	Responder   Contact
}

type findNodeRequest struct {
	target			*Contact
	receiver		*Contact
	responseChannel	chan findNodeResponse
}

type findNodeResponse struct {
	sender		*Contact
	contacts	[]Contact
}

// NewNetwork Constructor function for Network class
func NewNetwork(rt *RoutingTable, ht *ValueHashtable) *Network {
	network := &Network{}
	network.rt = rt
	network.ht = ht
	network.findNodeChannel = make(chan findNodeRequest)
	network.storeChannel = make(chan storeRequest)
	network.findDataChannel = make(chan findDataRequest)
	network.exitChannel = make(chan bool)
	return network
}

func (network *Network) handleChannels() {
	exit := false
	for !exit {
		select {
		case findNodeRequest := <- network.findNodeChannel:
			network.SendFindContactMessage(findNodeRequest)
		case storeRequest := <- network.storeChannel:
			network.SendStoreMessage(storeRequest)
		case findDataRequest := <- network.findDataChannel:
			network.SendFindDataMessage(findDataRequest)
		case exit = <- network.exitChannel:
		}
	}
}

// Helper function for listen
func (network *Network) handleListen(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var m msg
	err := decoder.Decode(&m)
	if err != nil {
		log.Println(err)
	}
	var mes string
	var cl []Contact = nil
	var d []byte = nil
	// depending on what message we got we run different
	switch m.Message {
	case "ping":
		// ping handle
		mes = "Response from ping"
	case "findcontact":
		// find contact handle
		mes = "Response from findcontact"
		cl = network.rt.FindClosestContacts(m.Target.ID, k)
	case "finddata":
		// find data handle
		// If a corresponding value is present on the recipient, the associated data is returned
		d = network.ht.Get(m.Hash)
		if d == nil {
			// Otherwise the RPC is equivalent to a FIND_NODE and a set of k triples is returned
			cl = network.rt.FindClosestContacts(m.Target.ID, k)
		}
		mes = "Response from finddata"
	case "store":
		// store handle
		network.ht.Put(m.Hash, m.Data)
		mes = "Response from store"
	default:
		log.Println("server received an invalid message")
		mes = "Response invalid message"
	}

	rm := response_msg{
		Message:     mes,
		ContactList: cl,
		Data:        d,
		Responder:   *network.rt.me,
	}

	r, err := json.Marshal(rm)
	if err != nil {
		log.Print(err)
	}
	// adds RPC sender to list
	go network.NetAddCont(m.Sender)

	fmt.Fprintf(rw, string(r))
}

func (network *Network) sendhelper(mes string, hash [HashSize]byte, data []byte, target *Contact, address string) response_msg {
	tm := msg{
		Message: mes,
		Hash:    hash,
		Data:    data,
		Sender:  *network.rt.me,
	}
	if target != nil {
		tm.Target = *target
	}
	requestBody, err := json.Marshal(tm)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(address, "message", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
		// maybe ping fail should be here
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshals
	var rm = response_msg{
		Message:     "error",
		ContactList: nil,
		Data:        nil,
		Responder:   Contact{},
	}
	err1 := json.Unmarshal(body, &rm)
	if err1 != nil {
		log.Println(err1)
	}
	return rm
}

// I guess you need to run this function as a go func
func (network *Network) Listen(address string, serveMux *http.ServeMux) {
	fmt.Println("Server starting on:", address)
	serveMux.HandleFunc("/msg", network.handleListen)
	log.Fatal(http.ListenAndServe(address, serveMux))
}

func (network *Network) SendPingMessage(receiver *Contact) bool {
	c1 := make(chan response_msg, 1)
	c2 := make(chan response_msg, 1)
	go func() {
		var nilHash [HashSize]byte
		rm := network.sendhelper("ping", nilHash, nil, nil, "http://"+receiver.Address+"/msg")
		c1 <- rm
		c2 <- rm
	}()

	if network.VibeCheck(c1) {
		rm := <-c2
		go network.NetAddCont(rm.Responder)
		return true
	}
	return false
}

func (network *Network) SendFindContactMessage(request findNodeRequest) {
	c1 := make(chan response_msg, 1)
	c2 := make(chan response_msg, 1)
	go func() {
		var nilHash [HashSize]byte
		rm := network.sendhelper("findcontact", nilHash, nil, request.target, "http://" + request.receiver.Address + "/msg")
		c1 <- rm
		c2 <- rm
	}()

	if network.VibeCheck(c1) {
		rm := <-c2
		go network.NetAddCont(rm.Responder)
		if rm.ContactList == nil {
			log.Println("Error: node has no contacts and returns nil")
		}
		request.responseChannel <- findNodeResponse{rm.Responder, rm.ContactList}
	}
}

// Retrieves the data from the receiver node using the hash key
func (network *Network) SendFindDataMessage(request findDataRequest) {
	c1 := make(chan response_msg, 1)
	c2 := make(chan response_msg, 1)
	go func() {
		rm := network.sendhelper("finddata", request.hash, nil, request.target, "http://" + request.receiver.Address + "/msg")
		c1 <- rm
		c2 <- rm
	}()

	if network.VibeCheck(c1) {
		rm := <-c2
		go network.NetAddCont(rm.Responder)
		if rm.ContactList == nil {
			log.Println("Error: node has no contacts and returns nil")
		}
		request.responseChannel <- findDataResponse{rm.Responder, rm.Data, rm.ContactList}
	}
}

// Tells the receiving node to store the data
func (network *Network) SendStoreMessage(request storeRequest) {
	network.sendhelper("store", request.hash, request.data, nil, "http://" + request.receiver.Address + "/msg")
}

func (network *Network) VibeCheck(c1 chan response_msg) bool {
	select {
	case res := <-c1:
		// Succeeds to get a response message
		if res.Message != "error" || res.Message == "" {
			return true
		}
		return false
	case <-time.After(3 * time.Second):
		// Times out
		fmt.Println("out of time, node is dead")
		return false
	}
}

func (network *Network) NetAddCont(contact Contact) {
	// if bucket is full
	if network.rt.buckets[network.rt.getBucketIndex(contact.ID)].Len() == k {
		// get last in list
		var last = network.rt.buckets[network.rt.getBucketIndex(contact.ID)].GetLast()
		// if it's not alive then we add, else we don't
		if !network.SendPingMessage(&last) {
			if !contact.ID.Equals(network.rt.me.ID) {
				network.rt.AddContact(contact)
			}
		} else {
			network.rt.AddContact(last)
		}
	} else {
		if !contact.ID.Equals(network.rt.me.ID) {
			network.rt.AddContact(contact)
		}
	}
}
