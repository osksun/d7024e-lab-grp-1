package main 

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv" )
 
var ipList [50]string 
var port string = "8080" 
func main() {
	myip()
	for i := 1; i < 51; i++ {
		ipList[i-1] = "172.17.0." + strconv.FormatInt(int64(i+1), 10)
	}
	fmt.Println(ipList)
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter text: ")
			text, _ := reader.ReadString('\n')
			id, err := strconv.ParseInt(text[:len(text)-1], 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			sendtoserver(id)
		}
	}()
	hostserver()
}
func sendtoserver(id int64) {
	fmt.Println("Sending get to container " + strconv.FormatInt(int64(id), 10))
	resp, err := http.Get("http://" + ipList[id-1] + ":" + port)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
}
func hostserver() {
	fmt.Println("Started server ...")
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":"+port, nil)
}
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a response from container with ip: %s", myip())
}
func myip() net.IP {
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}
	return nil
}
