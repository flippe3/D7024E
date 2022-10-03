package main

import (
	"testing"
)

func TestNetworkGetMyIp(t *testing.T) {
	recievedIp := GetMyIP()
	if recievedIp == "127.26.0.1" {
		t.Errorf("Less is very broken it says %s", recievedIp)
	}
}

func TestNetworkListenHandleConnection(t *testing.T) {
	// listener, err := net.Listen("tcp", ":81")
	// if err != nil {
	// 	t.Errorf("Can't start a tcp listener on port 81: %s", err)
	// }
	// conn, err := listener.Accept()
	// if conn
	// kademlia.network.HandleConnection()
	//kademlia.network.Listen(kademlia.dataMap)

	//kademlia2 := Kademlia{}
	//kademlia2.Join()
}

func TestNetworkCheckAliveAddContact(t *testing.T) {
	// listener, err := net.Listen("tcp", ":81")
	// if err != nil {
	// 	t.Errorf("Can't start a tcp listener on port 81: %s", err)
	// }
	// conn, err := listener.Accept()
	// if conn
	// kademlia.network.HandleConnection()
	//kademlia.network.Listen(kademlia.dataMap)

	//kademlia2 := Kademlia{}
	//kademlia2.Join()
}
