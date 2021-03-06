package d7024e

import (
	"strconv"
)

// Contacts command for the cli
func Contacts() Cmd{
	return Cmd{
		triggers: []string{"contacts", "c"},
		description: "Shows the contacts of the node",
		usage: "\"contacts\", \"c\"",
		action: func(cli *Cli, args ...string) string {
			buckets := cli.node.rt.buckets
			contactListStr := "\t| My address\t\t| My ID\n\t| " + cli.node.contact.Address + "\t| " + cli.node.contact.ID.String() + "\nContact List\n"
			contactListStr += "B_ID\t| Address\t\t| K_ID\t\t\t\t\t\t| Distance\n"
			for i := 0; i < len(buckets); i++ {
				for contact := buckets[i].list.Front(); contact != nil; contact = contact.Next() {
					contact := contact.Value.(Contact)
					contact.CalcDistance(cli.node.rt.me.ID)
					contactListStr +=  strconv.Itoa(i) + "\t| " + contact.Address + "\t| " + contact.ID.String() + "\t| " + contact.Distance.String() + "\n"
				}
			}
			return contactListStr
		},
	}
}