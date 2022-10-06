package main

import (
	"testing"
)

func TestBucketAddContains(t *testing.T) {
	var tempContact = NewContact(NewRandomKademliaID(), "172.0.1.4")
	var bucket = newBucket()
	bucket.AddContact(tempContact)
	bucket.AddContact(tempContact)
	if bucket.Contains(tempContact) == false {
		t.Errorf("bucket.Contains does not work")
	}
	var tempContact2 = NewContact(NewRandomKademliaID(), "172.0.1.5")
	if bucket.Contains(tempContact2) == true {
		t.Errorf("bucket.Contains does not work")
	}
}
