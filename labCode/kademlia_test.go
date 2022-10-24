package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestKademliaJoin(t *testing.T) {
	kademlia := Kademlia{}
	kademlia.Join()
	kademlia.Store("lol", 999999999)
	kademlia.LookupData("403926033d001b5279df37cbbe5287b7c7c267fa")
	kademlia.LookupData("403926033d001b5279df37cbbe5287b7c7c267ff")
	ch := make(chan int)
	go RunKademlia(&kademlia, ch)
	ch <- 0
	go CliParser(&kademlia, ch)
}

func TestKademliaFillBuckets(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	kademlia := Kademlia{network: &Network{routingTable: NewRoutingTable(NewContact(NewRandomKademliaID(), ""))}}
	rand.Seed(time.Now().UnixNano())
	newContact := NewContact(NewRandomKademliaID(), "")
	// Generate new kademlia IDs until it doesn't have bucket index 0
	for kademlia.network.routingTable.getBucketIndex(newContact.ID) == 0 {
		rand.Seed(time.Now().UnixNano())
		newContact = NewContact(NewRandomKademliaID(), "")
	}
	bucketIndex := kademlia.network.routingTable.getBucketIndex(newContact.ID)
	kademlia.network.routingTable.buckets[bucketIndex].AddContact(newContact)
	idList := kademlia.FillBuckets()
	if len(idList) != bucketIndex {
		t.Errorf("Lowest non-empty bucket index is %v, %v ids were generated.", bucketIndex-1, len(idList))
	}
	fmt.Println("myID: ", kademlia.network.routingTable.me.ID.String())
	fmt.Println("randomID: ", newContact.ID)
	for i, id := range idList {
		fmt.Println("bucket index: ", kademlia.network.routingTable.getBucketIndex(&id), " gets filled with ID: ", id.String())
		if kademlia.network.routingTable.getBucketIndex(&id) != len(idList)-i-1 {
			t.Errorf("id %v with bucket index %v was expected to have bucket index %v",
				id.String(), kademlia.network.routingTable.getBucketIndex(&id), len(idList)-i)
		}
	}
}

func TestKademliaHandleResponse(t *testing.T) {
	receivedCandidates := []Contact{NewContact(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), ""),
		NewContact(NewKademliaID("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), ""),
		NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "")}
	contactShortlist := ContactCandidates{contacts: []Contact{NewContact(NewKademliaID("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), ""),
		NewContact(NewKademliaID("cccccccccccccccccccccccccccccccccccccccc"), "")}}
	queriedContacts := ContactCandidates{}
	queriedContact := contactShortlist.contacts[1]
	target := NewKademliaID("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	contactShortlist.contacts[0].CalcDistance(target)
	contactShortlist.contacts[1].CalcDistance(target)
	kademlia := Kademlia{network: &Network{routingTable: NewRoutingTable(NewContact(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), ""))}}
	kademlia.HandleResponse(&queriedContacts, queriedContact, receivedCandidates, &contactShortlist, target)

	if contactShortlist.Len() != 3 {
		t.Errorf("contactShortlist did not have expected length 3, instead has length %v", contactShortlist.Len())
	}
	if queriedContacts.Len() != 1 {
		t.Errorf("queriedContacts did not have expected length 1, instead has length %v", queriedContacts.Len())
	}
	if !contactShortlist.contacts[0].ID.Equals(NewKademliaID("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")) {
		t.Errorf("index 0 of contactShortlist did not have expected kademliaID aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa, instead got kademlia ID %s", contactShortlist.contacts[0].ID.String())
	}
	if !contactShortlist.contacts[1].ID.Equals(NewKademliaID("cccccccccccccccccccccccccccccccccccccccc")) {
		t.Errorf("index 1 of contactShortlist did not have expected kademliaID cccccccccccccccccccccccccccccccccccccccc, instead got kademlia ID %s", contactShortlist.contacts[1].ID.String())
	}
	if !contactShortlist.contacts[2].ID.Equals(NewKademliaID("0000000000000000000000000000000000000000")) {
		t.Errorf("index 2 of contactShortlist did not have expected kademliaID 0000000000000000000000000000000000000000, instead got kademlia ID %s", contactShortlist.contacts[2].ID.String())
	}
	if !queriedContacts.contacts[0].ID.Equals(NewKademliaID("cccccccccccccccccccccccccccccccccccccccc")) {
		t.Errorf("index 0 of queriedContacts did not have expected kademliaID cccccccccccccccccccccccccccccccccccccccc, instead got kademlia ID %s", queriedContacts.contacts[0].ID.String())
	}
}
