package main

import (
	"testing"
)

func TestFindClosestContacts(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2fffffff00000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("f111111400000000000000000000000000000000"), 4)

	if !contacts[0].ID.Equals(NewKademliaID("2fffffff00000000000000000000000000000000")) {
		t.Errorf("Did not find closest contact. contacts[0].ID = %s", contacts[0].ID)
	}
	for i := 0; i < len(contacts)-1; i++ {
		contacts[i].CalcDistance(NewKademliaID("2111111400000000000000000000000000000000"))
		contacts[i+1].CalcDistance(NewKademliaID("2111111400000000000000000000000000000000"))
		if !contacts[i].Less(&contacts[i+1]) {
			t.Errorf("contact[%v] was not closer than contact[%v]", i, i+1)
		}
	}
}

func TestGetBucketIndex(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"))
	if rt.getBucketIndex(NewKademliaID("ffffffff00000000000000000000000000000001")) != 159 {
		t.Errorf("Couldn't index closest contact correctly.")
	}
	if rt.getBucketIndex(NewKademliaID("00000000ffffffffffffffffffffffffffffffff")) != 0 {
		t.Errorf("Couldn't index furthest away contact correctly.")
	}
	if rt.getBucketIndex(NewKademliaID("ff7fffff00000000000000000000000000000000")) != 8 {
		t.Errorf("Expected kademlia ID ff7fffff00000000000000000000000000000000 to have index 8. Calculated index %v", rt.getBucketIndex(NewKademliaID("ff7fffff00000000000000000000000000000000")))
	}
}

func TestLowestNonEmptyBucketIndex(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000"))
	rt.AddContact(NewContact(NewKademliaID("ffffffff00000000000000000000000000000001"), ""))
	if rt.LowestNonEmptyBucketIndex() != 159 {
		t.Errorf("Lowest non-empty bucket index: %v, expected 159", rt.LowestNonEmptyBucketIndex())
	}
	rt.AddContact(NewContact(NewKademliaID("ff7fffff00000000000000000000000000000000"), ""))
	if rt.LowestNonEmptyBucketIndex() != 8 {
		t.Errorf("Lowest non-empty bucket index: %v, expected 8", rt.LowestNonEmptyBucketIndex())
	}
	rt.AddContact(NewContact(NewKademliaID("00000000ffffffffffffffffffffffffffffffff"), ""))
	if rt.LowestNonEmptyBucketIndex() != 0 {
		t.Errorf("Lowest non-empty bucket index: %v, expected 0", rt.LowestNonEmptyBucketIndex())
	}
}
