import "net"
package d7024e type Network struct {
}
func Listen(ip string, port int) {
	ln, err:= net.Listen("tcp", ip+":"+port)
	if err != nil {
	// error occured 
	}
	for {
		conn, err :=ln.Accept()		
		if err != nil {
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}
func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}
func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}
func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

