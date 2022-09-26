package main

import (
	"crypto/sha1"
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
func (network *Network) Listen(dataMap map[string]string) {
	listener, err := net.Listen("tcp", ":81")
	if err != nil {
		fmt.Println("Error occured while trying to listen: ", err)
	} else {
		fmt.Println("Listening on ", listener.Addr().String())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error occured while trying to accept connection: ", err)
		} else {
			network.HandleConnection(conn, dataMap)
		}
		conn.Close()
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
		strMsg := network.routingTable.me.ID.String() + "," + network.routingTable.me.Address
		fmt.Println(strMsg)
		conn.Write([]byte(strMsg))

	} else if string([]rune(inp)[0]) == "p" {
		// SendPing 🏓
		// Add id and ip to routing table
		// returns my id and ip COMMA SEPERATED :3
		var inp = strings.Split(string(b), ",")
		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
	} else if string([]rune(inp)[0]) == "c" {
		// FindContact 👷‍♀️
		// returns that contacts nearest 4 contacts from my routing table
		// Expects inp[3] to be targetID
		var inp = strings.Split(string(b), ",")
		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
		getClosestFour := network.routingTable.FindClosestContacts(NewKademliaID(inp[3]), 4)

		// Create the message
		var strMsg string
		for i, k := range getClosestFour {
			contactArr := strings.Split(k.String(), "\"")
			strMsg += contactArr[1] + "," + contactArr[3]
			if i != len(getClosestFour)-1 {
				strMsg += ","
			}
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
		// TODO: This needs to send back something saying its found or not found
		// the data.
		if ok {
			conn.Write([]byte(dataMap[data]))
		} else {
			getClosestFour := network.routingTable.FindClosestContacts(NewKademliaID(inp[3]), 4)

			strMsg := network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + ","
			for i, k := range getClosestFour {
				strMsg += k.String()
				if i != len(getClosestFour)-1 {
					strMsg += ","
				}
			}
			conn.Write([]byte(strMsg))
		}

	} else if string([]rune(inp)[0]) == "s" {
		// StoreMessage 🚛
		var inp = strings.Split(string(b), ",")
		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))

		h := sha1.New()
		hash := string(h.Sum([]byte(inp[3])))
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
		strMessage := "p," + network.routingTable.me.ID.String() + "," + network.routingTable.me.Address
		conn.Write([]byte(strMessage))
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
		messageStr := "c," + network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + "," + target.String()
		_, writeErr := conn.Write([]byte(messageStr))
		if writeErr != nil {
			return nil, writeErr
		} else {
			b := make([]byte, 255)
			conn.Read(b)
			conn.Close()
			bArr := strings.Split(string(b), ",")
			contacts := []Contact{}
			for i := 0; i < len(bArr); i = i + 2 {
				contacts = append(contacts, NewContact(NewKademliaID(bArr[i]), bArr[i+1]))
			}
			for _, contact := range contacts {
				if contact.ID.Equals(network.routingTable.me.ID) {
					continue // Ignore self
				}
				network.CheckAliveAddContact(contact) // Update routing table
			}
			fmt.Println("SendFindContactMessage to ", contact.String(), "with target: ", target.String(), ", found", len(contacts), "contacts: ", ContactsString(contacts))
			return contacts, nil
		}
	}
}

// d (Has to be comma seperated)
func (network *Network) SendFindDataMessage(contact *Contact, hash string) {
	// TODO
}

// s
func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
