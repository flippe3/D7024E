package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Sets up a reader that splits the input into 2 strings
// [0] = operation, [1] = value
func CliParser(kademlia *Kademlia, exit chan int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		read, _ := reader.ReadString('\n')
		fmt.Println(CliHandler(strings.Fields(read), kademlia, exit))
	}
}

func CliHandler(s []string, kademlia *Kademlia, exit chan int) string {
	if len(s) == 0 {
		return ""
	}
	switch operation := s[0]; operation {
	case "put":
		if len(s) != 3 {
			return "Expected exactly 2 arguments for command 'put'"
		}
		r := sha1.Sum([]byte(s[1]))
		ttl, err := strconv.Atoi(s[2])
		if err != nil {
			return "Expected second argument to be a ttl integer"
		}
		return "put " + s[1] + " " + s[2] + " sent store RPCs with hash: " + hex.EncodeToString(r[:]) + ", to contacts: " + ContactsString(kademlia.Store(s[1], ttl))
	case "get":
		if len(s) != 2 {
			return "Expected exactly 1 argument for command 'get'"
		} else if len(s[1]) != 40 {
			return "Argument provided to get is not a SHA-1 hash"
		}
		data, contacts := kademlia.LookupData(s[1])
		if data != "" {
			return "get " + s[1] + " found the data: " + data
		}
		return "get " + s[1] + " did not find the data.\n 'get' found the contactsFound contacts: " + ContactsString(contacts)
	case "exit":
		if len(s) != 1 {
			return "Expected 0 arguments for command 'exit'"
		} else {
			exit <- 0
			return ""
		}
	default:
		return "operation: " + operation + " not found"
	}
}
