package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const k int = 4

type Kademlia struct {
	network *Network
	dataMap map[string]string
}

// Functions in this file are iterative

// j
func (kademlia *Kademlia) Join() {
	// Setup your kademlia
	myIP := GetMyIP()
	myKademliaID := NewRandomKademliaID()
	fmt.Println("Generating new kademlia ID", myKademliaID, "for ip: ", myIP)
	routingTable := NewRoutingTable(NewContact(myKademliaID, myIP))
	kademlia.network = &Network{routingTable: routingTable}

	slice := strings.Split(myIP, ".")
	for i := 2; i <= 255; i++ {
		addr := slice[0] + "." + slice[1] + "." + slice[2] + "." + strconv.Itoa(i) + ":81"
		if addr == myIP+":81" {
			continue
		}
		conn, connErr := net.DialTimeout("tcp", addr, 1e6) // TTL: 1 ms
		if connErr != nil {
			continue
		} else {
			_, writeErr := conn.Write([]byte("j"))
			if writeErr != nil {
				fmt.Println("Failed to send join message", writeErr)
			}

			b := make([]byte, 255)
			_, readErr := conn.Read(b)
			if readErr != nil {
				fmt.Println("2Read error", readErr)
			}
			conn.Close()
			bArr := strings.Split(string(b), ",")

			bootstrappingContact := NewContact(NewKademliaID(bArr[0]), bArr[1])
			fmt.Println("Adding bootstrapping contact: ", bootstrappingContact.String())
			routingTable.AddContact(NewContact(NewKademliaID(bArr[0]), bArr[1]))
			contacts := kademlia.LookupContact(routingTable.me.ID)
			fmt.Println("Looked up following contacts when looking for self: ", ContactsString(contacts))
			lowestNonEmptyIndex := kademlia.network.routingTable.LowestNonEmptyBucketIndex(contacts)
			fmt.Println("Calculated the following kademliaIDs to look up when filling buckets: ")
			idList := kademlia.FillBuckets(lowestNonEmptyIndex)
			for _, id := range idList {
				fmt.Println(id.String())
			}
			for _, id := range idList {
				kademlia.LookupContact(&id)
			}
			return
		}
	}
}

// -------------------- CREATE UNIT TEST --------------------
// Fills buckets that are of higher index than the lowest non empty
func (kademlia *Kademlia) FillBuckets(lowestNonEmptyIndex int) []KademliaID {
	var prefixOnes []byte
	var invertedBit []byte
	var suffixOnes []byte
	for i := 0; i < IDLength-1; i++ {
		prefixOnes = append(prefixOnes, byte(255))
		invertedBit = append(invertedBit, byte(0))
		suffixOnes = append(suffixOnes, byte(0))
	}
	prefixOnes = append(prefixOnes, byte(254))
	invertedBit = append(invertedBit, byte(2))
	suffixOnes = append(suffixOnes, byte(1))
	index := IDLength*8 - lowestNonEmptyIndex
	for i := 0; i < index-1; i++ {
		ShiftLeft(prefixOnes)
		ShiftLeft(invertedBit)
		ShiftLeft(suffixOnes)
		suffixOnes[IDLength-1] = suffixOnes[IDLength-1] | byte(1)
	}
	myID, _ := hex.DecodeString(kademlia.network.routingTable.me.ID.String())
	randomID, _ := hex.DecodeString(NewRandomKademliaID().String())
	// Establish the full kademliaIDs and look them up
	lookupIDs := []KademliaID{}
	for prefixOnes[0] != byte(0) {
		lookupID := Or(Xor(And(prefixOnes, myID), invertedBit), And(randomID, suffixOnes))
		lookupIDs = append(lookupIDs, (*NewKademliaID(hex.EncodeToString(lookupID))))
		ShiftLeft(prefixOnes)
		ShiftLeft(invertedBit)
		ShiftLeft(suffixOnes)
		suffixOnes[IDLength-1] = suffixOnes[IDLength-1] | byte(1)
	}
	return lookupIDs
}

// -------------------- CREATE UNIT TEST --------------------
// Shifts the given byte slice left by one
func ShiftLeft(data []byte) {
	for i := 0; i < len(data)-1; i++ {
		data[i] = data[i]<<1 | data[i+1]>>7
	}
	data[len(data)-1] <<= 1
}

/*
// Shifts the given byte slice right by one
func ShiftRight(data []byte) {
	for i := len(data) - 1; i > 0; i-- {
		data[i] = data[i]>>1 | data[i-1]<<7
	}
	data[0] >>= 1
}
*/

// -------------------- CREATE UNIT TEST --------------------
// And operator for byte slices
func And(data1 []byte, data2 []byte) []byte {
	var data []byte
	for i := 0; i < len(data1); i++ {
		data = append(data, data1[i]&data2[i])
	}
	return data
}

// -------------------- CREATE UNIT TEST --------------------
// Or operator for byte slices
func Or(data1 []byte, data2 []byte) []byte {
	var data []byte
	for i := 0; i < len(data1); i++ {
		data = append(data, data1[i]|data2[i])
	}
	return data
}

// -------------------- CREATE UNIT TEST --------------------
// Xor operator for byte slices
func Xor(data1 []byte, data2 []byte) []byte {
	var data []byte
	for i := 0; i < len(data1); i++ {
		data = append(data, data1[i]^data2[i])
	}
	return data
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	probedContacts := ContactCandidates{contacts: []Contact{}}
	unprobedContacts := ContactCandidates{contacts: []Contact{}}
	closest := Contact{}
	for probedContacts.Len() < k {
		if unprobedContacts.IsEmpty() {
			// Get new contact candidate from own routing table
			candidates := kademlia.network.routingTable.FindClosestContacts(target, 999999999999999)
			for i := 0; i < len(candidates); i++ {
				if probedContacts.Contains(&candidates[i]) {
					continue
				} else {
					fmt.Println("Found unprobed contact: ", candidates[i].String())
					unprobedContacts.Append(candidates[i : i+1])
					if closest.ID == nil {
						closest = candidates[i : i+1][0]
					}
					break
				}
			}
			if unprobedContacts.IsEmpty() {
				// No new contact candidates found, returns less than k found contacts
				fmt.Println("Did not find any new contacts to probe")
				probedContacts.Sort()
				return probedContacts.GetContacts(k)
			} else {
				continue
			}
		} else {
			fmt.Println("Probing contact: ", unprobedContacts.contacts[0].String())
			// Contact next unprobed contact
			receivedCandidates, err := kademlia.network.SendFindContactMessage(&unprobedContacts.contacts[0], target)
			if err != nil {
				// No response
				fmt.Println("1No response from ", unprobedContacts.contacts[0].String())
				unprobedContacts = *unprobedContacts.Remove(0)
				continue
			}
			unprobedContacts = kademlia.HandleFindContactResponse(receivedCandidates, &unprobedContacts, &probedContacts, target)
			unprobedContacts.Sort()
			fmt.Println("Receieved contactcandidates: ", ContactsString(receivedCandidates), "unprobedContacts: ", unprobedContacts.String(), "probedContacts: ", probedContacts.String())
			if unprobedContacts.IsEmpty() || closest.Less(&unprobedContacts.GetContacts(1)[0]) {
				// Did not find a closer contact
				fmt.Println("Did not find a closer contact")
				for i := 0; i < k && !unprobedContacts.IsEmpty(); i++ {
					receivedCandidates, err := kademlia.network.SendFindContactMessage(&unprobedContacts.contacts[0], target)
					if err != nil {
						// No response
						fmt.Println("2No response from ", unprobedContacts.contacts[0].String())
						unprobedContacts = *unprobedContacts.Remove(0)
						continue
					}
					unprobedContacts = kademlia.HandleFindContactResponse(receivedCandidates, &unprobedContacts, &probedContacts, target)
					unprobedContacts.Sort()
					fmt.Println("Receieved contactcandidates: ", ContactsString(receivedCandidates), "unprobedContacts: ", unprobedContacts.String(), "probedContacts: ", probedContacts.String())
				}
				probedContacts.Sort()
				return probedContacts.GetContacts(k)
			} else {
				// Found a closer contact
				fmt.Println("Found closer contact: ", unprobedContacts.GetContacts(1)[0].String())
				closest = unprobedContacts.GetContacts(1)[0]
			}
		}
	}
	// Successfully found k probed closest contacts
	probedContacts.Sort()
	return probedContacts.GetContacts(k)
}

// -------------------- CREATE UNIT TEST --------------------
// LookupContact helper
// Returns new list of unprobedContacts (since remove is used)
func (kademlia *Kademlia) HandleFindContactResponse(
	receivedCandidates []Contact, unprobedContacts *ContactCandidates, probedContacts *ContactCandidates, target *KademliaID) ContactCandidates {
	probedContacts.Append(unprobedContacts.contacts[0:1])
	unprobedContacts = unprobedContacts.Remove(0)
	for i := 0; i < len(receivedCandidates); i++ {
		if receivedCandidates[i].ID.Equals(kademlia.network.routingTable.me.ID) ||
			probedContacts.Contains(&receivedCandidates[i]) ||
			unprobedContacts.Contains(&receivedCandidates[i]) {
			continue // Ignore self, already contacted or already appended
		} else {
			receivedCandidates[i].CalcDistance(target)
			unprobedContacts.Append(receivedCandidates[i : i+1])
		}
	}
	return *unprobedContacts
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
