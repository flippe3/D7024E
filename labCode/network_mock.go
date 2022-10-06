package main

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func (network *Network) MockHandleConnection(inp string, dataMap map[string]string) string {
	if string([]rune(inp)[0]) == "j" {
		writing := network.routingTable.me.ID.String() + "," + network.routingTable.me.Address + ","
		return writing
	} else if string([]rune(inp)[0]) == "c" {
		var inp = strings.Split(inp, ",")
		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
		getClosestFour := network.routingTable.FindClosestContacts(NewKademliaID(inp[3]), 4) // Expects inp[3] to be targetID

		writing := ""
		for _, k := range getClosestFour {
			contactArr := strings.Split(k.String(), "\"")
			writing += contactArr[1] + "," + contactArr[3] + ","
		}
		return writing
	} else if string([]rune(inp)[0]) == "d" {
		var inp = strings.Split(inp, ",")
		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
		data, ok := dataMap[inp[3]]
		if ok {
			writing := data + ","
			return writing
		} else {
			getClosestFour := network.routingTable.FindClosestContacts(NewKademliaID(inp[3]), 4)
			writing := ""
			for _, k := range getClosestFour {
				contactArr := strings.Split(k.String(), "\"")
				writing += contactArr[1] + "," + contactArr[3] + ","
			}
			return writing
		}
	} else if string([]rune(inp)[0]) == "s" {
		var inp = strings.Split(inp, ",")
		network.CheckAliveAddContact(NewContact(NewKademliaID(inp[1]), inp[2]))
		r := sha1.Sum([]byte(inp[3]))
		hash := hex.EncodeToString(r[:])
		dataMap[hash] = inp[3]
		return ""
	}
	return ""
}
