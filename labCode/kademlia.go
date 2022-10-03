package main

import (
	"crypto/sha1"
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
	for i := 1; i <= 255; i++ {
		addr := slice[0] + "." + slice[1] + "." + slice[2] + "." + strconv.Itoa(i) + ":81"
		if addr == myIP+":81" {
			continue
		}
		conn, connErr := net.DialTimeout("tcp", addr, 1e6) // TTL: 1 ms
		if connErr != nil {
			continue
		} else {
			_, writeErr := conn.Write([]byte("j,"))
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
			fmt.Println("Looked up following contacts when looking for self: ", ContactsString(kademlia.LookupContact(routingTable.me.ID)))
			fmt.Println("Calculated the following kademliaIDs to look up when filling buckets: ")
			idList := kademlia.FillBuckets()
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

// Fills buckets that are of higher index than the lowest non empty
func (kademlia *Kademlia) FillBuckets() []KademliaID {
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
	index := IDLength*8 - kademlia.network.routingTable.LowestNonEmptyBucketIndex()
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

// And operator for byte slices
func And(data1 []byte, data2 []byte) []byte {
	var data []byte
	for i := 0; i < len(data1); i++ {
		data = append(data, data1[i]&data2[i])
	}
	return data
}

// Or operator for byte slices
func Or(data1 []byte, data2 []byte) []byte {
	var data []byte
	for i := 0; i < len(data1); i++ {
		data = append(data, data1[i]|data2[i])
	}
	return data
}

// Xor operator for byte slices
func Xor(data1 []byte, data2 []byte) []byte {
	var data []byte
	for i := 0; i < len(data1); i++ {
		data = append(data, data1[i]^data2[i])
	}
	return data
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	queriedContacts := ContactCandidates{contacts: []Contact{}}
	contactShortlist := ContactCandidates{contacts: kademlia.network.routingTable.FindClosestContacts(target, k)}
	fmt.Println("Closest contacts in own routing table: ", ContactsString(contactShortlist.GetContacts(k)))
	for i := 0; i < k; i++ {
		if contactShortlist.Len() > i && !queriedContacts.Contains(&contactShortlist.contacts[i]) {
			fmt.Println("Probing contact: ", contactShortlist.contacts[i].String())
			receivedCandidates, err := kademlia.network.SendFindContactMessage(&contactShortlist.contacts[i], target)
			if err != nil {
				// No response
				fmt.Println("1No response from ", contactShortlist.contacts[i].String())
				contactShortlist = *contactShortlist.Remove(i)
				i--
				continue
			}
			queriedContacts.Append([]Contact{contactShortlist.contacts[i]})
			closestBefore := contactShortlist.contacts[0]
			kademlia.AppendToShortlist(receivedCandidates, &contactShortlist, target)
			fmt.Println("Received contactcandidates: ", ContactsString(receivedCandidates), "\ncontactShortlist: ", contactShortlist.String(), "\nqueriedContacts: ", queriedContacts.String())
			i = -1
			if closestBefore.ID.Equals(contactShortlist.contacts[0].ID) {
				// Did not find a closer contact
				fmt.Println("Did not find a closer contact")
				// Probe k closest not already probed
				for j := 0; j < k && contactShortlist.Len() > j; j++ {
					if queriedContacts.Contains(&contactShortlist.contacts[j]) {
						continue
					}
					if kademlia.network.SendPingMessage(&contactShortlist.contacts[j]) {
						queriedContacts.Append([]Contact{contactShortlist.contacts[j]})
					} else {
						// No response
						fmt.Println("2No response from ", contactShortlist.contacts[j].String())
						contactShortlist = *contactShortlist.Remove(j)
						j--
					}
				}
				return contactShortlist.GetContacts(k)
			}
		}
	}
	return contactShortlist.GetContacts(k)
}

// LookupContact helper
func (kademlia *Kademlia) AppendToShortlist(
	receivedCandidates []Contact, contactShortlist *ContactCandidates, target *KademliaID) {
	for i := 0; i < len(receivedCandidates); i++ {
		if receivedCandidates[i].ID.Equals(kademlia.network.routingTable.me.ID) || contactShortlist.Contains(&receivedCandidates[i]) {
			continue // Ignore self and already known about nodes
		} else {
			receivedCandidates[i].CalcDistance(target)
			contactShortlist.Append(receivedCandidates[i : i+1])
		}
	}
	contactShortlist.Sort()
}

func (kademlia *Kademlia) LookupData(hash string) (string, []Contact) {
	queriedContacts := ContactCandidates{contacts: []Contact{}}
	contactShortlist := ContactCandidates{contacts: kademlia.network.routingTable.FindClosestContacts(NewKademliaID(hash), k)}
	fmt.Println("Closest contacts in own routing table: ", ContactsString(contactShortlist.GetContacts(k)))
	for i := 0; i < k; i++ {
		if contactShortlist.Len() > i && !queriedContacts.Contains(&contactShortlist.contacts[i]) {
			fmt.Println("Probing contact: ", contactShortlist.contacts[i].String())
			data, receivedCandidates, err := kademlia.network.SendFindDataMessage(&contactShortlist.contacts[i], hash)
			if err != nil {
				// No response
				fmt.Println("8No response from ", contactShortlist.contacts[i].String())
				contactShortlist = *contactShortlist.Remove(i)
				i--
				continue
			}
			if data != "" {
				// Store data in closest queried contact which did not return the data
				if queriedContacts.Len() != 0 {
					queriedContacts.Sort()
					kademlia.network.SendStoreMessage(&queriedContacts.contacts[0], data)
				}
				return data, nil
			}
			queriedContacts.Append([]Contact{contactShortlist.contacts[i]})
			closestBefore := contactShortlist.contacts[0]
			kademlia.AppendToShortlist(receivedCandidates, &contactShortlist, NewKademliaID(hash))
			fmt.Println("Received contactcandidates: ", ContactsString(receivedCandidates), "\ncontactShortlist: ", contactShortlist.String(), "\nqueriedContacts: ", queriedContacts.String())
			i = -1
			if closestBefore.ID.Equals(contactShortlist.contacts[0].ID) {
				// Did not find a closer contact
				fmt.Println("Did not find a closer contact")
				// Probe k closest not already probed
				for j := 0; j < k && contactShortlist.Len() > j; j++ {
					if queriedContacts.Contains(&contactShortlist.contacts[j]) {
						continue
					}
					data, _, err := kademlia.network.SendFindDataMessage(&contactShortlist.contacts[j], hash)
					if err != nil {
						// No response
						fmt.Println("9No response from ", contactShortlist.contacts[j].String())
						contactShortlist = *contactShortlist.Remove(j)
						j--
						continue
					}
					if data != "" {
						// Store data in closest queried contact which did not return the data
						if queriedContacts.Len() != 0 {
							queriedContacts.Sort()
							kademlia.network.SendStoreMessage(&queriedContacts.contacts[0], data)
						}
						return data, nil
					}
					queriedContacts.Append([]Contact{contactShortlist.contacts[j]})
				}
				return "", contactShortlist.GetContacts(k)
			}
		}
	}
	return "", contactShortlist.GetContacts(k)
}

func (kademlia *Kademlia) Store(data string) []Contact {
	var id KademliaID = sha1.Sum([]byte(data))
	fmt.Println("Looking to store data: ", data, " with hash: ", id.String())
	contacts := kademlia.LookupContact(&id)
	for _, contact := range contacts {
		kademlia.network.SendStoreMessage(&contact, data)
	}
	return contacts
}
