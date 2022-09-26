package main

import (
	"fmt"
	"sort"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}

// NewContact returns a new instance of a Contact
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

// String returns a simple string representation of a Contact
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

// ContactCandidates definition
// stores an array of Contacts
type ContactCandidates struct {
	contacts []Contact
}

// Append an array of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

// GetContacts returns the first count number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	if count > candidates.Len() {
		count = candidates.Len()
	}
	return candidates.contacts[:count]
}

// Sort the Contacts in ContactCandidates
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

func (candidates *ContactCandidates) IsEmpty() bool {
	return candidates.Len() == 0
}

func (candidates *ContactCandidates) Contains(contact *Contact) bool {
	for i := 0; i < candidates.Len(); i++ {
		if candidates.contacts[i].ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}

// Removes contact at index
func (candidates *ContactCandidates) Remove(index int) *ContactCandidates {
	var newCandidates []Contact
	newCandidates = append(newCandidates, candidates.GetContacts(index)...)
	newCandidates = append(newCandidates, candidates.GetContacts(candidates.Len())[index+1:]...)
	return &ContactCandidates{contacts: newCandidates}
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}

func (candidates ContactCandidates) String() string {
	return ContactsString(candidates.GetContacts(candidates.Len()))
}

func ContactsString(contacts []Contact) string {
	str := ""
	for i, contact := range contacts {
		if i == len(contacts)-1 {
			str += contact.String()
		} else {
			str += contact.String() + ","
		}
	}
	return str
}
