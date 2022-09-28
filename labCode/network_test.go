package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNetworkGetMyIp(t *testing.T) {
	recievedIp := GetMyIP()
	if recievedIp == "127.26.0.1" {
		t.Errorf("Less is very broken it says %s", recievedIp)
	}
}

func TestNetworkHandleConnection(t *testing.T) {
	kademlia := Kademlia{}
	fmt.Print(kademlia)
	kademlia.Join()
	// Make sure bucket with index 0 is empty at least once
	var i int64 = 0
	for ; kademlia.network.routingTable.LowestNonEmptyBucketIndex() == 0; i++ {
		rand.Seed(i)
		kademlia.Join()
	}
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
