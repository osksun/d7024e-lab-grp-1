package d7024e

import (
)

// Contacts command for the cli
func Contacts() Cmd{
	return Cmd{
		triggers: []string{"contacts", "c"},
		description: "Shows the contacts of the node",
		usage: "\"contacts\", \"c\"",
		action: func(cli *Cli, args ...string) string {
			var contactList []Contact
			buckets := cli.node.rt.buckets
			for i := 0; i < len(buckets); i++ {
				for contact := buckets[i].list.Front(); contact != nil; contact = contact.Next() {
					contactList = append(contactList, contact.Value.(Contact))
				}
			}
			contactListStr := "Address\t|\tKademliaID\n"
			for i := 0; i < len(contactList); i++ {
				contactListStr += contactList[i].Address + "\t|\t" + contactList[i].ID.String() + "\n"
			}
			return contactListStr
		},
	}
}