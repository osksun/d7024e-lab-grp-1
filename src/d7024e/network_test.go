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
		t.Errorf("The network is not of type network")
	}
	fmt.Printf("TestNewNetwork finished running with status OK\n")
}

type TestHttpHandler struct {
	n Network
}

func (h *TestHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.n.handleListen(w, r)
}
func TestHandleListen(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	nct := NewContact(kID1, "localhost:8001")
	rt := NewRoutingTable(nct)
	ht := NewValueHashtable()
	n := NewNetwork(rt, ht)
	h := &TestHttpHandler{n: *n}

	server := httptest.NewServer(h)
	defer server.Close()
	fmt.Println(server.URL)
	// Make a test request
	resp, err := http.POST(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}

// // Helper function for listen
// func (network *Network) handleListen(rw http.ResponseWriter, req *http.Request) {
// 	decoder := json.NewDecoder(req.Body)
// 	var m msg
// 	err := decoder.Decode(&m)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	var mes string
// 	var cl []Contact = nil
// 	var d []byte = nil
// 	// depending on what message we got we run different
// 	switch m.Message {
// 	case "ping":
// 		// ping handle
// 		mes = "Response from ping"
// 	case "findcontact":
// 		// find contact handle
// 		mes = "findcontact response"
// 		cl = network.rt.FindClosestContacts(m.Target.ID, IDLength) // K = 20 here
// 	case "finddata":
// 		// find data handle
// 		d = network.ht.Get(m.Hash)
// 		mes = "Response from finddata"
// 	case "store":
// 		// store handle
// 		// PUT NEEDS A STRING KEY ASSOCIATED WITH THE DATA
// 		network.ht.Put(m.Hash, m.Data)
// 		mes = "Response from store"
// 	default:
// 		log.Println("server received an invalid message")
// 		mes = "Response: invalid message"
// 	}

// 	rm := response_msg{
// 		Message:     mes,
// 		ContactList: cl,
// 		Data:        d,
// 		Responder:   *network.rt.me,
// 	}

// 	r, err := json.Marshal(rm)
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	// adds RPC sender to list
// 	go network.NetAddCont(m.Sender)

// 	fmt.Fprintf(rw, string(r))
// }

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
