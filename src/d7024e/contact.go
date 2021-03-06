package d7024e

import (
	"fmt"
	"sort"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	ID       *KademliaID
	Address  string
	Distance *KademliaID
	queried 	bool
}

// NewContact returns a new instance of a Contact
func NewContact(id *KademliaID, address string) *Contact {
	return &Contact{id, address, nil, false}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.Distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.Distance.Less(otherContact.Distance)
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

// NewContactCandidates constructor for creating an instance of ContactCandidates
func NewContactCandidates(contacts []Contact) *ContactCandidates {
	return &ContactCandidates{contacts}
}

// Append an array of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

// GetContacts returns the first count number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
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

// RemoveDuplicates removes duplicate contacts based on KademliaID
func (candidates *ContactCandidates) RemoveDuplicates() {
	keys := make(map[KademliaID]bool)
	list := []Contact{}
	for i := 0; i < candidates.Len(); i++ {
		if _, value := keys[*candidates.contacts[i].ID]; !value {
			keys[*candidates.contacts[i].ID] = true
			list = append(list, candidates.contacts[i])
		}
	}
	candidates.contacts = list
}

// remove removes a contact by given *KademliaID, the order of the list will not be maintained
func (candidates *ContactCandidates) removeContacts(contacts []*Contact) {
	for i := 0; i < len(contacts); i++ {
		candidates.remove(contacts[i].ID)
	}
}

// remove removes a contact by given *KademliaID, the order of the list will not be maintained
func (candidates *ContactCandidates) remove(kademliaID *KademliaID) {
	for i := 0; i < candidates.Len(); i++ {
		if candidates.contacts[i].ID.Equals(kademliaID) {
			candidates.contacts [i] = candidates.contacts[len(candidates.contacts ) - 1]
			candidates.contacts = candidates.contacts[:len(candidates.contacts ) - 1]
		}
	}
}

// removeContact removes a contact by given *KademliaID, the order of the list will not be maintained
func removeContact(contacts []*Contact, kademliaID *KademliaID) []*Contact{
	for i := 0; i < len(contacts); i++ {
		if contacts[i].ID.Equals(kademliaID) {
			contacts[i] = contacts[len(contacts) - 1]
			return contacts[:len(contacts) - 1]
		}
	}
	return contacts
}
