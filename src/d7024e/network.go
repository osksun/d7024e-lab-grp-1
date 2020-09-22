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
	rt *RoutingTable
	ht *ValueHashtable
}

type msg struct {
	Message string
	Hash    string
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

// NewNetwork Constructor function for Network class
func NewNetwork(rt *RoutingTable, ht *ValueHashtable) *Network {
	network := &Network{}
	network.rt = rt
	network.ht = ht
	return network
}

// Helper function for listen
func (network *Network) handleListen(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var m msg
	err := decoder.Decode(&m)
	if err != nil {
		log.Println(err)
	}
	log.Println(m.Message)

	var mes string
	var cl []Contact = nil
	var d []byte = nil
	// depending on what message we got we run different
	switch m.Message {
	case "ping":
		// ping handle
		log.Println("server ping")
		mes = "Response from ping"
	case "findcontact":
		// find contact handle
		log.Println("server findcontact")
		mes = "findcontact response"
		cl = network.rt.FindClosestContacts(m.Target.ID, 20) // K = 20 here
	case "finddata":
		// find data handle
		log.Println("server finddata")
		d = network.ht.Get(m.Hash)
		mes = "Response from finddata"
	case "store":
		// store handle
		log.Println("server store")
		// PUT NEEDS A STRING KEY ASSOCIATED WITH THE DATA
		network.ht.Put("keyHERE", m.Data)
		mes = "Response from store"
	default:
		log.Println("server received an invalid message")
		mes = "Response: invalid message"
	}

	rm := response_msg{
		Message:     mes,
		ContactList: cl,
		Data:        d,
		Responder:   *network.rt.me,
	}

	r, err := json.Marshal(rm)

	// adds RPC sender to list
	network.NetAddCont(m.Sender)

	fmt.Fprintf(rw, string(r))
}

func (network *Network) sendhelper(mes string, hash string, data []byte, target *Contact, address string) response_msg {
	tm := msg{
		Message: mes,
		Hash:    hash,
		Data:    data,
		Target:  *target,
		Sender:  *network.rt.me,
	}
	requestBody, err := json.Marshal(tm)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("sending ... ")
	resp, err := http.Post("http://"+address+"/msg", "message", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
		// maybe ping fail should be here
	}
	log.Println(resp)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshals
	log.Println(string(body))
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
	log.Println("here is the rm response:")
	log.Println(rm)
	return rm
}

// I guess you need to run this function as a go func
func (network *Network) Listen(address string, serveMux *http.ServeMux) {
	fmt.Println("Server starting on:", address)
	serveMux.HandleFunc("/msg", network.handleListen)
	log.Fatal(http.ListenAndServe(address, serveMux))
}

func (network *Network) SendPingMessage(receiver *Contact) bool {
	// TODO
	c1 := make(chan response_msg, 1)
	c2 := make(chan response_msg, 1)
	go func() {
		rm := network.sendhelper("ping", "", nil, nil, receiver.Address)
		c1 <- rm
		c2 <- rm
	}()

	if network.VibeCheck(c1) {
		//rm := <-c2
		//network.NetAddCont(rm.Responder)
		return true
	}
	return false
}

func (network *Network) SendFindContactMessage(target *Contact, receiver *Contact) []Contact {
	// TODO
	c1 := make(chan response_msg, 1)
	c2 := make(chan response_msg, 1)
	go func() {
		rm := network.sendhelper("findcontact", "", nil, target, receiver.Address)
		c1 <- rm
		c2 <- rm
	}()

	if network.VibeCheck(c1) {
		rm := <-c2
		network.NetAddCont(rm.Responder)
		return rm.ContactList
	}
	return nil
}

// Retrieves the data from the receiver node using the hash key
func (network *Network) SendFindDataMessage(receiver *Contact, hash string) {
	// TODO
	rm := network.sendhelper("finddata", hash, nil, nil, receiver.Address)
	log.Println(rm.Message)
}

// Tells the receiving node to store the data
func (network *Network) SendStoreMessage(receiver *Contact, data []byte) {
	// TODO
	rm := network.sendhelper("store", "", data, nil, receiver.Address)
	log.Println(rm.Message)
}

func (network *Network) VibeCheck(c1 chan response_msg) bool {
	select {
	case res := <-c1:
		// Succeeds to get a response message
		if res.Message != "error" {
			log.Println("Succeeds the vibecheck")
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
	if network.rt.buckets[network.rt.getBucketIndex(contact.ID)].Len() == bucketSize {
		// get last in list
		var last = network.rt.buckets[network.rt.getBucketIndex(contact.ID)].GetLast()
		// if it's not alive then we add, else we don't
		if !network.SendPingMessage(&last) {
			network.rt.AddContact(contact)
		}
	} else {
		network.rt.AddContact(contact)
	}
}
