package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"strings"
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
func (network *Network) Listen(dataMap map[string]string, chConn chan net.Conn) {
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

func (network *Network) HandleConnection(conn net.Conn, dataMap map[string]string) {
	b := make([]byte, 255)
	n, err := conn.Read(b)
	if err != nil {
		fmt.Println("1Read error", err)
	}
	fmt.Println("Inside network: ", "n: ", n, "b: ", b)

	// j = Join = 106
	// p = SendPing = 112
	// c = FindContact 😍 = 99
	// d = FindData = 100
	// s = StoreMessage = 115

	var inp = string(b)
	if string([]rune(inp)[0]) == "j" {
		// Join 🤗
		// Sends back my id and ip COMMA SEPERATED :)
		strMsg := network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + ","
		fmt.Println(strMsg)
		conn.Write([]byte(strMsg))

	} else if string([]rune(inp)[0]) == "p" {
		// SendPing 🏓
		// Add id and ip to routing table
		// returns my id and ip COMMA SEPERATED :3

		// Commenting out these 2 lines below because they lead to infinite pings
		//var inp = strings.Split(string(b), ",")
		//network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
	} else if string([]rune(inp)[0]) == "c" {
		// FindContact 👷‍♀️
		// returns that contacts nearest 4 contacts from my routing table
		// Expects inp[3] to be targetID
		var inp = strings.Split(string(b), ",")

		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))

		getClosestFour := network.routingTable.FindClosestContacts(NewKademliaID(inp[3]), 4)
		// Create the message
		var strMsg string
		for _, k := range getClosestFour {
			contactArr := strings.Split(k.String(), "\"")
			strMsg += contactArr[1] + "," + contactArr[3] + ","
		}
		fmt.Println("Looking for closest contacts to ID: ", NewKademliaID(inp[3]), ", found : ", strMsg)
		conn.Write([]byte(strMsg))
	} else if string([]rune(inp)[0]) == "d" {
		// FindDataMessage 📀
		// Look for the data in dataMap
		// ex: nodeid,nodeip,hash
		var inp = strings.Split(string(b), ",")

		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
		data, ok := dataMap[inp[3]]
		if ok {
			fmt.Println("Found data: ", data)
			conn.Write([]byte(data + ","))
		} else {
			getClosestFour := network.routingTable.FindClosestContacts(NewKademliaID(inp[3]), 4)
			// Create the message
			var strMsg string
			for _, k := range getClosestFour {
				contactArr := strings.Split(k.String(), "\"")
				strMsg += contactArr[1] + "," + contactArr[3] + ","
			}
			fmt.Println("Looking for closest contacts to ID: ", NewKademliaID(inp[3]), ", found : ", strMsg)
			conn.Write([]byte(strMsg))
		}
	} else if string([]rune(inp)[0]) == "s" {
		// StoreMessage 🚛
		var inp = strings.Split(string(b), ",")

		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
		r := sha1.Sum([]byte(inp[3]))
		hash := hex.EncodeToString(r[:])
		fmt.Println("Storing data ", inp[3], " with hash: ", hash)
		dataMap[hash] = inp[3]
	}
}

// Checks if contact should be added to the bucket
func (network *Network) CheckAliveAddContact(contact Contact) {
	fmt.Println("CheckAliveAddContact: ", contact.String())
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

		// Commenting out these 2 lines below because they lead to infinite pings
		//strMessage := "p," + network.routingTable.me.ID.String() + "," + network.routingTable.me.Address
		//conn.Write([]byte(strMessage))

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
		messageStr := "c," + network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + "," + target.String() + ","
		_, writeErr := conn.Write([]byte(messageStr))
		if writeErr != nil {
			return nil, writeErr
		} else {
			b := make([]byte, 255)
			conn.Read(b)
			conn.Close()
			bArr := strings.Split(string(b), ",")
			contacts := []Contact{}
			for i := 0; i < len(bArr)-1; i = i + 2 {
				contacts = append(contacts, NewContact(NewKademliaID(bArr[i]), bArr[i+1]))
			}
			fmt.Println("SendFindContactMessage to ", contact.String(), "with target: ", target.String(), ", found", len(contacts), "contacts: ", ContactsString(contacts))
			return contacts, nil
		}
	}
}

// d (Has to be comma seperated) (string, []Contact, error)
func (network *Network) SendFindDataMessage(contact *Contact, hash string) (string, []Contact, error) {
	conn, connErr := net.DialTimeout("tcp", contact.Address+":81", 1e6) // TTL: 1 ms
	if connErr != nil {
		return "", nil, fmt.Errorf("4No response from " + contact.String())
	} else {
		network.CheckAliveAddContact(*contact) // Update routing table
		messageStr := "d," + network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + "," + hash + ","
		_, writeErr := conn.Write([]byte(messageStr))
		if writeErr != nil {
			return "", nil, writeErr
		} else {
			b := make([]byte, 255)
			conn.Read(b)
			conn.Close()
			bArr := strings.Split(string(b), ",")
			if len(bArr) == 2 {
				fmt.Println("SendFindDataMessage to ", contact.String(), "with target: ", hash, ", found the data ", bArr[0])
				return bArr[0], nil, nil
			} else {
				contacts := []Contact{}
				for i := 0; i < len(bArr)-1; i = i + 2 {
					contacts = append(contacts, NewContact(NewKademliaID(bArr[i]), bArr[i+1]))
				}
				fmt.Println("SendFindDataMessage to ", contact.String(), "with target: ", hash, ", found", len(contacts), "contacts: ", ContactsString(contacts))
				return "", contacts, nil
			}
		}
	}
}

// s
func (network *Network) SendStoreMessage(contact *Contact, data string) {
	conn, connErr := net.DialTimeout("tcp", contact.Address+":81", 1e6) // TTL: 1 ms
	if connErr != nil {
		fmt.Println("5No response from " + contact.String())
	} else {
		network.CheckAliveAddContact(*contact) // Update routing table
		messageStr := "s," + network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + "," + data + ","
		_, writeErr := conn.Write([]byte(messageStr))
		if writeErr != nil {
			fmt.Println("WriteError to "+contact.String(), " while trying to write ", data)
		}
	}
}
