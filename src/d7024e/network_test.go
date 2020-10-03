package d7024e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
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
			t.Errorf("Sender got '%s' want '%s'", rm.Message, expected_response[i])
		}

	}
	fmt.Println("TestHandleSender finished running with status OK")
}

// Not sure how to test this one

// // I guess you need to run this function as a go func
// func (network *Network) Listen(address string, serveMux *http.ServeMux) {
// 	fmt.Println("Server starting on:", address)
// 	serveMux.HandleFunc("/msg", network.handleListen)
// 	log.Fatal(http.ListenAndServe(address, serveMux))
// }

func TestSendPing(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	nct := NewContact(kID1, "localhost:8001")
	rt := NewRoutingTable(nct)
	ht := NewValueHashtable()
	n := NewNetwork(rt, ht)

	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct2 := NewContact(kID2, "localhost:8002")
	rt2 := NewRoutingTable(nct2)
	ht2 := NewValueHashtable()
	n2 := NewNetwork(rt2, ht2)

	go n2.Listen(n2.rt.me.Address, http.NewServeMux())
	time.Sleep(1 * time.Second)
	resp := n.SendPingMessage(nct2)
	if resp != true {
		t.Error("SendPingMessage failed the test.")
	}
	fmt.Println("TestSendPing finished running with status OK")
}

func TestSendFindContactMessage(t *testing.T) {
	/*
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	nct := NewContact(kID1, "localhost:8001")
	rt := NewRoutingTable(nct)
	ht := NewValueHashtable()
	n := NewNetwork(rt, ht)
	*/
	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct2 := NewContact(kID2, "localhost:8003")
	rt2 := NewRoutingTable(nct2)
	ht2 := NewValueHashtable()
	n2 := NewNetwork(rt2, ht2)

	kID3, _ := NewKademliaID("FFFFFF0000000000000000000000000000000000")
	nct3 := NewContact(kID3, "localhost:8004")

	n2.rt.AddContact(*nct3)

	go n2.Listen(n2.rt.me.Address, http.NewServeMux())
	time.Sleep(1 * time.Second)
	/*
	resp := n.SendFindContactMessage(nct3, nct2)
	if len(resp) != 1 {
		t.Error("SendFindContactMessage failed the test.")
	}
	*/
	fmt.Println("TestSendFindContactMessage finished running with status OK")
}

func TestSendFindDataMessage(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	nct := NewContact(kID1, "localhost:8001")
	rt := NewRoutingTable(nct)
	ht := NewValueHashtable()
	n := NewNetwork(rt, ht)

	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct2 := NewContact(kID2, "localhost:8004")
	rt2 := NewRoutingTable(nct2)
	ht2 := NewValueHashtable()
	n2 := NewNetwork(rt2, ht2)

	go n2.Listen(n2.rt.me.Address, http.NewServeMux())
	time.Sleep(1 * time.Second)

	data := []byte("data")
	key := "key"
	hash := Hash([]byte(key))

	n2.ht.Put(hash, data)

	resp := n.SendFindDataMessage(hash, nct2)
	if reflect.DeepEqual(resp, data) == false {
		t.Error("SendFindDataMessage failed the test.")
	}
	fmt.Println("TestSendFindDataMessage finished running with status OK")
}

func TestSendStoreMessage(t *testing.T) {
	kID1, _ := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	nct := NewContact(kID1, "localhost:8001")
	rt := NewRoutingTable(nct)
	ht := NewValueHashtable()
	n := NewNetwork(rt, ht)

	kID2, _ := NewKademliaID("FFFFFFF000000000000000000000000000000000")
	nct2 := NewContact(kID2, "localhost:8005")
	rt2 := NewRoutingTable(nct2)
	ht2 := NewValueHashtable()
	n2 := NewNetwork(rt2, ht2)

	go n2.Listen(n2.rt.me.Address, http.NewServeMux())
	time.Sleep(1 * time.Second)

	data := []byte("data")
	key := "key"
	hash := Hash([]byte(key))

	n.SendStoreMessage(nct2, hash, data)
	if reflect.DeepEqual(n2.ht.Get(hash), data) == false {
		t.Error("SendStoreMessage failed the test.")
	}
	fmt.Println("TestSendStoreMessage finished running with status OK")
}

func TestVibeCheck(t *testing.T) {

}

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
