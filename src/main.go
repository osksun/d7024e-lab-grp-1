package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"./d7024e"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	joinip := getJoin()
	log.Println("Joining to node: " + joinip)
	node := d7024e.NewNode(checkip() + ":8000", "")
	node.SpinupNode(false, true, joinip)

}

func getJoin() string {
	ip := checkip()
	fmt.Println("My IP: " + ip)
	joinip := ""
	if ip != "172.17.0.2" {
		max, _ := extract(ip)
		min := 2
		lastnum := rand.Intn(max-min) + min
		joinip = "172.17.0." + strconv.Itoa(lastnum)
	}
	return joinip
}

func checkip() string {
	//flag := false
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				return ipnet.IP.String()

				//os.Stdout.WriteString("IP: " + ipnet.IP.String() + "\n")
			}
		}
	}
	return "nothing"
}

func extract(ip string) (int, error) {
	str1 := ip

	re := regexp.MustCompile(`([0-9]+)`)

	//fmt.Printf("Pattern: %v\n", re.String()) // print pattern
	fmt.Println(re.MatchString(str1)) // true

	submatchall := re.FindAllString(str1, -1)
	return strconv.Atoi(submatchall[len(submatchall)-1])
}
