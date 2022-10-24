package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Network struct {
	routingTable *RoutingTable
}

// The functions in this file are primitive

// Returns your ip address with respect to the correct subnet mask (16)
func GetMyIP() (address string) {
	nwiList, _ := net.Interfaces()
	for _, nwi := range nwiList {
		addrs, _ := nwi.Addrs()
		for _, addr := range addrs {
			matched, _ := regexp.MatchString(".*/16", addr.String())
			if matched {
				address, _, _ = strings.Cut(addr.String(), "/")
			}
		}
	}
	return address
}

// Listens for all incoming traffic on port 81
func (network *Network) Listen(dataMap map[string]Object, chConn chan net.Conn) {
	listener, err := net.Listen("tcp", ":81")
	if err != nil {
		fmt.Println("Error occured while trying to listen: ", err)
	} else {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error occured while trying to accept connection: ", err)
			} else {
				chConn <- conn
			}
		}
	}
}

func (network *Network) HandleConnection(conn net.Conn, dataMap map[string]Object) {
	b := make([]byte, 255)
	conn.Read(b)
	inp := string(b)
	if string([]rune(inp)[0]) == "j" { // Join 🤗
		conn.Write([]byte(network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + ","))
	} else if string([]rune(inp)[0]) == "c" { // FindContact 👷‍♀️
		var inp = strings.Split(string(b), ",")
		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
		getClosestFour := network.routingTable.FindClosestContacts(NewKademliaID(inp[3]), 4) // Expects inp[3] to be targetID

		var strMsg string
		for _, k := range getClosestFour {
			contactArr := strings.Split(k.String(), "\"")
			strMsg += contactArr[1] + "," + contactArr[3] + ","
		}
		conn.Write([]byte(strMsg))
	} else if string([]rune(inp)[0]) == "d" { // FindDataMessage 📀
		var inp = strings.Split(string(b), ",")
		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
		data, ok := dataMap[inp[3]]
		if ok {
			if data.storageTime.Add(time.Duration(data.ttl * 1e9)).After(time.Now()) {
				dataMap[inp[3]] = Object{data: data.data, ttl: data.ttl, storageTime: time.Now()}
				conn.Write([]byte("d," + data.data + "," + strconv.Itoa(data.ttl) + ","))
				return
			}
			delete(dataMap, inp[3])
		}
		getClosestFour := network.routingTable.FindClosestContacts(NewKademliaID(inp[3]), 4)
		var strMsg string
		for _, k := range getClosestFour {
			contactArr := strings.Split(k.String(), "\"")
			strMsg += contactArr[1] + "," + contactArr[3] + ","
		}
		conn.Write([]byte(strMsg))
	} else if string([]rune(inp)[0]) == "s" { // StoreMessage 🚛
		var inp = strings.Split(string(b), ",")
		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
		r := sha1.Sum([]byte(inp[3]))
		hash := hex.EncodeToString(r[:])
		ttl, _ := strconv.Atoi(inp[4])
		dataMap[hash] = Object{data: inp[3], ttl: ttl, storageTime: time.Now()}
	}
}

// Checks if contact should be added to the bucket
func (network *Network) CheckAliveAddContact(contact Contact) {
	//fmt.Println("CheckAliveAddContact: ", contact.String())
	bucketIndex := network.routingTable.getBucketIndex(contact.ID)
	bucket := network.routingTable.buckets[bucketIndex]
	if bucket.Contains(contact) {
		bucket.AddContact(contact)
	} else {
		if bucket.Len() == bucketSize {
			leastRecentlySeen := bucket.list.Back().Value.(Contact)
			if !network.SendPingMessage(&leastRecentlySeen) {
				bucket.list.Remove(bucket.list.Back())
				bucket.AddContact(contact)
			}
		} else {
			bucket.AddContact(contact)
		}
	}
}

// p
// Returns true if a response is received within the TTL (1 ms)
func (network *Network) SendPingMessage(contact *Contact) (response bool) {
	conn, connErr := net.DialTimeout("tcp", contact.Address+":81", 1e6) // TTL: 1 ms
	if connErr != nil {
		return false
	} else {
		network.CheckAliveAddContact(*contact) // Update routing table
		conn.Close()
		return true
	}
}

// c
// Returns sorted list of contacts
func (network *Network) SendFindContactMessage(contact *Contact, target *KademliaID) ([]Contact, error) {
	conn, connErr := net.DialTimeout("tcp", contact.Address+":81", 1e6) // TTL: 1 ms
	if connErr != nil {
		return nil, fmt.Errorf("3No response from " + contact.String())
	} else {
		network.CheckAliveAddContact(*contact) // Update routing table
		conn.Write([]byte("c," + network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + "," + target.String() + ","))
		b := make([]byte, 255)
		conn.Read(b)
		conn.Close()
		bArr := strings.Split(string(b), ",")
		contacts := []Contact{}
		for i := 0; i < len(bArr)-1; i = i + 2 {
			contacts = append(contacts, NewContact(NewKademliaID(bArr[i]), bArr[i+1]))
		}
		//fmt.Println("SendFindContactMessage to ", contact.String(), "with target: ", target.String(), ", found", len(contacts), "contacts: ", ContactsString(contacts))
		return contacts, nil
	}
}

// d (Has to be comma seperated) (string, []Contact, error)
func (network *Network) SendFindDataMessage(contact *Contact, hash string) (Object, []Contact, error) {
	conn, connErr := net.DialTimeout("tcp", contact.Address+":81", 1e6) // TTL: 1 ms
	if connErr != nil {
		return Object{data: ""}, nil, fmt.Errorf("4No response from " + contact.String())
	} else {
		network.CheckAliveAddContact(*contact) // Update routing table
		conn.Write([]byte("d," + network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + "," + hash + ","))
		b := make([]byte, 255)
		conn.Read(b)
		conn.Close()
		bArr := strings.Split(string(b), ",")
		if bArr[0] == "d" {
			//fmt.Println("SendFindDataMessage to ", contact.String(), "with target: ", hash, ", found the data ", bArr[1])
			ttl, _ := strconv.Atoi((bArr[2]))
			return Object{data: bArr[1], ttl: ttl}, nil, nil
		}
		contacts := []Contact{}
		for i := 0; i < len(bArr)-1; i = i + 2 {
			contacts = append(contacts, NewContact(NewKademliaID(bArr[i]), bArr[i+1]))
		}
		//fmt.Println("SendFindDataMessage to ", contact.String(), "with target: ", hash, ", found", len(contacts), "contacts: ", ContactsString(contacts))
		return Object{data: ""}, contacts, nil
	}
}

// s
func (network *Network) SendStoreMessage(contact *Contact, data string, ttl int) {
	conn, connErr := net.DialTimeout("tcp", contact.Address+":81", 1e6) // TTL: 1 ms
	if connErr != nil {
		//fmt.Println("5No response from " + contact.String())
	} else {
		network.CheckAliveAddContact(*contact) // Update routing table
		conn.Write([]byte("s," + network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + "," + data + "," + strconv.Itoa(ttl) + ","))
		conn.Close()
	}
}
