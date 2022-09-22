package main

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type Network struct {
	myIP         string
	routingTable *RoutingTable
}

// The functions in this file are primitive

func Listen(ip string, port int) {
	listener, err := net.Listen("tcp", ":81")
	if err != nil {
		fmt.Print("We are not listening", err)
	} else {
		fmt.Print("Listening on ", listener)
	}
	conn, err := listener.Accept()
	fmt.Print("postlisten")
	if err != nil {
		fmt.Print("We are not accepted", err)
	}
	var b []byte
	n, err := conn.Read(b)
	if err != nil {
		fmt.Print("Read error", err)
	}
	fmt.Print(n, b)
	conn.Close()
}

func Join() *Network {
	myIP := GetMyIP()
	slice := strings.Split(myIP, ".")
	var addr string
	var conn net.Conn
	var connErr error
	var writeErr error
	for i := 0; i <= 4; i++ {
		addr = slice[0] + "." + slice[1] + "." + slice[2] + "." + strconv.Itoa(i) + ":81"
		if addr == myIP+":81" {
			continue
		}
		conn, connErr = net.DialTimeout("tcp", addr, 1e9)
		if connErr != nil {
			fmt.Print(connErr)
			continue
		}
		_, writeErr = conn.Write([]byte{72, 101, 74})
		if writeErr != nil {
			fmt.Print("writeError inside JOIN", writeErr)
		}
		conn.Close()
		//fmt.Print("n INSIDE JoIN: ", n)
	}
	return &Network{myIP: myIP, routingTable: NewRoutingTable(NewContact(NewRandomKademliaID(), myIP))}
}

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

// Port nr 82
func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

// Port nr 83
func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

// Port nr 84
func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

// Port nr 85
func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
