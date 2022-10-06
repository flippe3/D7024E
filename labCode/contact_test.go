package main

import (
	"testing"
)

func TestContactString(t *testing.T) {
	var id = NewRandomKademliaID()
	var c = NewContact(id, "127.0.1.4")
	if c.String() != "contact(\""+id.String()+"\", \"127.0.1.4\")" {
		t.Errorf("Contact.String() should be %s but is %s", "contact(\""+id.String()+"\", \"127.0.1.4\")", c.String())
	}
}

func TestContactIsEmpty(t *testing.T) {
	var candidates = ContactCandidates{contacts: []Contact{}}
	if candidates.IsEmpty() == false {
		t.Errorf("candidates.IsEmpty() should be empty but is not.")
	}
	var id = NewRandomKademliaID()
	var c = NewContact(id, "172.0.1.4")
	candidates.Append([]Contact{c})
	if candidates.IsEmpty() == true {
		t.Errorf("candidates.IsEmpty() should be non-empty but is empty.")
	}
}

func TestContactContains(t *testing.T) {
	var candidates = ContactCandidates{contacts: []Contact{}}
	var c = NewContact(NewRandomKademliaID(), "172.0.1.4")
	var c2 = NewContact(NewRandomKademliaID(), "172.0.1.5")
	candidates.Append([]Contact{c})
	if candidates.Contains(&c) == false {
		t.Errorf("candidates.Contains() ERROR, should contain %s and does not", c.String())
	} else if candidates.Contains(&c2) == true {
		t.Errorf("candidates.Contains() ERROR, should not contain %s but does", c2.String())
	}
}

func TestContactRemove(t *testing.T) {
	var c = NewContact(NewRandomKademliaID(), "172.0.1.4")
	var c2 = NewContact(NewRandomKademliaID(), "172.0.1.5")
	var candidates = ContactCandidates{contacts: []Contact{c, c2}}
	var newcand = candidates.Remove(1)
	if newcand.Len() != 1 {
		t.Errorf("candidates.Remove() should make candidates length 1 but is length %d", candidates.Len())
	}
	// TODO: Check that it removed the correct contact
}

func TestCandidatesString(t *testing.T) {
	var c = NewContact(NewRandomKademliaID(), "172.0.1.4")
	var c2 = NewContact(NewRandomKademliaID(), "172.0.1.5")
	var candidates = ContactCandidates{contacts: []Contact{c, c2}}
	if candidates.String() != c.String()+","+c2.String() {
		t.Errorf("candidates.String() should make %s but makes %s", c.String()+","+c2.String(), candidates.String())
	}
}

func TestCandidatesGetContacts(t *testing.T) {
	contact := NewContact(NewRandomKademliaID(), "")
	contactCandidates := ContactCandidates{contacts: []Contact{contact}}
	if len(contactCandidates.GetContacts(999)) != 1 {
		t.Errorf("Did not have expected length 1")
	}
	if !contactCandidates.GetContacts(999)[0].ID.Equals(contact.ID) {
		t.Errorf("IDs did not equate")
	}
}
