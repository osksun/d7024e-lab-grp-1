package d7024e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// type Network struct {
// 	rt *RoutingTable
// 	ht *ValueHashtable
// }

func TestNewNetwork(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	nct := NewContact(kID1, "localhost:8001")
	rt := NewRoutingTable(nct)
	ht := NewValueHashtable()
	n := NewNetwork(rt, ht)
	nType := fmt.Sprintf("%T", n)
	if nType != "*d7024e.Network" {
		t.Error("The network is not of type network")
	}
	fmt.Println("TestNewNetwork finished running with status OK")
}

type TestHttpHandler struct {
	n Network
}

func (h *TestHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.n.handleListen(w, r)
}
func TestHandleSender(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	nct := NewContact(kID1, "localhost:8001")
	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct2 := NewContact(kID2, "localhost:8002")
	rt := NewRoutingTable(nct)
	ht := NewValueHashtable()
	n := NewNetwork(rt, ht)
	h := &TestHttpHandler{n: *n}

	server := httptest.NewServer(h)
	defer server.Close()

	hash := Hash([]byte("something"))
	data := []byte("data")
	target := *nct2

	test_messages := [5]string{"ping", "findcontact", "finddata", "store", "invalidmessage"}
	expected_response := [5]string{"Response from ping", "Response from findcontact", "Response from finddata", "Response from store", "Response invalid message"}

	for i := 0; i < len(test_messages); i++ {
		mes := test_messages[i]
		rm := n.sendhelper(mes, hash, data, &target, server.URL)

		if rm.Message != expected_response[i] {
			t.Errorf("Sender didn't get the expected response. %s", expected_response[i])
		}

	}
	fmt.Println("TestHandleSender finished running with status OK")
}

// func (network *Network) sendhelper(mes string, hash [HashSize]byte, data []byte, target *Contact, address string) response_msg {
// 	tm := msg{
// 		Message: mes,
// 		Hash:    hash,
// 		Data:    data,
// 		Sender:  *network.rt.me,
// 	}
// 	if target != nil {
// 		tm.Target = *target
// 	}
// 	requestBody, err := json.Marshal(tm)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	resp, err := http.Post("http://"+address+"/msg", "message", bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		log.Fatalln(err)
// 		// maybe ping fail should be here
// 	}

// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// Unmarshals
// 	var rm = response_msg{
// 		Message:     "error",
// 		ContactList: nil,
// 		Data:        nil,
// 		Responder:   Contact{},
// 	}
// 	err1 := json.Unmarshal(body, &rm)
// 	if err1 != nil {
// 		log.Println(err1)
// 	}
// 	return rm
// }

// // I guess you need to run this function as a go func
// func (network *Network) Listen(address string, serveMux *http.ServeMux) {
// 	fmt.Println("Server starting on:", address)
// 	serveMux.HandleFunc("/msg", network.handleListen)
// 	log.Fatal(http.ListenAndServe(address, serveMux))
// }

// func (network *Network) SendPingMessage(receiver *Contact) bool {
// 	// TODO
// 	c1 := make(chan response_msg, 1)
// 	c2 := make(chan response_msg, 1)
// 	go func() {
// 		var nilHash [20]byte
// 		rm := network.sendhelper("ping", nilHash, nil, nil, receiver.Address)
// 		c1 <- rm
// 		c2 <- rm
// 	}()

// 	if network.VibeCheck(c1) {
// 		//rm := <-c2
// 		//network.rt.AddContact(rm.Responder)
// 		return true
// 	}
// 	return false
// }

// func (network *Network) SendFindContactMessage(target *Contact, receiver *Contact) []Contact {
// 	// TODO
// 	c1 := make(chan response_msg, 1)
// 	c2 := make(chan response_msg, 1)
// 	go func() {
// 		var nilHash [20]byte
// 		rm := network.sendhelper("findcontact", nilHash, nil, target, receiver.Address)
// 		c1 <- rm
// 		c2 <- rm
// 	}()

// 	if network.VibeCheck(c1) {
// 		rm := <-c2
// 		go network.NetAddCont(rm.Responder)
// 		if rm.ContactList == nil {
// 			log.Println("Error: node has no contacts and returns nil")
// 		}
// 		return rm.ContactList
// 	}
// 	return nil
// }

// // Retrieves the data from the receiver node using the hash key
// func (network *Network) SendFindDataMessage(hash [HashSize]byte, receiver *Contact) []byte {
// 	rm := network.sendhelper("finddata", hash, nil, nil, receiver.Address)
// 	return rm.Data
// }

// // Tells the receiving node to store the data
// func (network *Network) SendStoreMessage(receiver *Contact, hash [HashSize]byte, data []byte) {
// 	network.sendhelper("store", hash, data, nil, receiver.Address)
// }

// func (network *Network) VibeCheck(c1 chan response_msg) bool {
// 	select {
// 	case res := <-c1:
// 		// Succeeds to get a response message
// 		if res.Message != "error" || res.Message == "" {
// 			return true
// 		}
// 		return false
// 	case <-time.After(3 * time.Second):
// 		// Times out
// 		fmt.Println("out of time, node is dead")
// 		return false
// 	}
// }

// func (network *Network) NetAddCont(contact Contact) {
// 	// if bucket is full
// 	if network.rt.buckets[network.rt.getBucketIndex(contact.ID)].Len() == IDLength {
// 		// get last in list
// 		var last = network.rt.buckets[network.rt.getBucketIndex(contact.ID)].GetLast()
// 		// if it's not alive then we add, else we don't
// 		if !network.SendPingMessage(&last) {
// 			if !contact.ID.Equals(network.rt.me.ID) {
// 				network.rt.AddContact(contact)
// 			}
// 		} else {
// 			network.rt.AddContact(last)
// 		}
// 	} else {
// 		if !contact.ID.Equals(network.rt.me.ID) {
// 			network.rt.AddContact(contact)
// 		}
// 	}
// }
