package main

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"
)

func TestNetworkGetMyIp(t *testing.T) {
	recievedIp := GetMyIP()
	if recievedIp == "127.26.0.1" {
		t.Errorf("Less is very broken it says %s", recievedIp)
	}
}

func TestNetworkHandleConnection(t *testing.T) {
	dataMap := make(map[string]string)
	network := Network{routingTable: NewRoutingTable(NewContact(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), "127.26.0.1"))}
	network.routingTable.AddContact(NewContact(NewKademliaID("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), "127.26.0.3"))
	network.routingTable.AddContact(NewContact(NewKademliaID("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), "127.26.0.4"))
	network.routingTable.AddContact(NewContact(NewKademliaID("cccccccccccccccccccccccccccccccccccccccc"), "127.26.0.5"))
	network.routingTable.AddContact(NewContact(NewKademliaID("dddddddddddddddddddddddddddddddddddddddd"), "127.26.0.6"))
	network.routingTable.AddContact(NewContact(NewKademliaID("eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"), "127.26.0.7"))
	r := network.MockHandleConnection("j", dataMap)
	if r != network.routingTable.me.ID.String()+","+network.routingTable.me.Address+"," {
		t.Errorf("Unexpected write while handling j, wrote %v", r)
	}
	r = network.MockHandleConnection("c,0000000000000000000000000000000000000000,127.26.0.2,cccccccccccccccccccccccccccccccccccccccc,", dataMap)
	if r != "cccccccccccccccccccccccccccccccccccccccc,127.26.0.5,dddddddddddddddddddddddddddddddddddddddd,127.26.0.6,eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee,127.26.0.7,aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa,127.26.0.3," {
		t.Errorf("Unexpected write while handling c, wrote %v", r)
	}
	if !network.routingTable.buckets[0].Contains(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "127.26.0.2")) {
		t.Errorf("Did not add new contact in routing table.")
	}
	r = network.MockHandleConnection("d,0000000000000000000000000000000000000000,127.26.0.2,aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d,", dataMap)
	if r != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa,127.26.0.3,bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb,127.26.0.4,cccccccccccccccccccccccccccccccccccccccc,127.26.0.5,dddddddddddddddddddddddddddddddddddddddd,127.26.0.6," {
		t.Errorf("Unexpected write while handling d, wrote %v", r)
	}
	dataMap["aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"] = "hello"
	r = network.MockHandleConnection("d,0000000000000000000000000000000000000000,127.26.0.2,aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d,", dataMap)
	if r != "hello," {
		t.Errorf("Unexpected write while handling d, wrote %v", r)
	}
	network.MockHandleConnection("s,0000000000000000000000000000000000000000,127.26.0.2,world,", dataMap)
	q := sha1.Sum([]byte("world"))
	hash := hex.EncodeToString(q[:])
	if dataMap[hash] != "world" {
		t.Errorf("Could not find stored data 'world' in datamap")
	}
	if network.MockHandleConnection("q", dataMap) != "" {
		t.Errorf("Did not do nothing")
	}
}

func TestNetworkCheckAliveAddContact(t *testing.T) {
	network := Network{routingTable: NewRoutingTable(NewContact(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), ""))}
	network.routingTable.AddContact(NewContact(NewKademliaID("0fffffffffffffffffffffffffffffffffffffff"), ""))
	network.CheckAliveAddContact(NewContact(NewKademliaID("0fffffffffffffffffffffffffffffffffffffff"), ""))
	network.CheckAliveAddContact(NewContact(NewKademliaID("00ffffffffffffffffffffffffffffffffffffff"), ""))
	if network.routingTable.buckets[0].list.Len() != 2 {
		t.Errorf("Bucket with index 0 did not have the expected length 2, instead has length %v", network.routingTable.buckets[0].list.Len())
	}
	if !network.routingTable.buckets[0].list.Front().Value.(Contact).ID.Equals(NewKademliaID("00ffffffffffffffffffffffffffffffffffffff")) {
		t.Errorf("Front of bucket is not most recently seen.")
	}
}
