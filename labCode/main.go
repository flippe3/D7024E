package main

import (
	"net"
)

func main() {
	kademlia := Kademlia{dataMap: make(map[string]string)}
	kademlia.Join()
	exit := make(chan int)
	go CliParser(&kademlia, exit)
	RunKademlia(&kademlia, exit)
}

func RunKademlia(kademlia *Kademlia, exit chan int) {
	chConn := make(chan net.Conn)
	go kademlia.network.Listen(kademlia.dataMap, chConn)
	for {
		select {
		case conn := <-chConn:
			kademlia.network.HandleConnection(conn, kademlia.dataMap)
			conn.Close()
		case <-exit:
			return
		}
	}
}
