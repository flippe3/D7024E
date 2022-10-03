package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestKademliaJoin(t *testing.T) {
	// kademlia := Kademlia{}
	// fmt.Print(kademlia)
	// kademlia.Join()
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
