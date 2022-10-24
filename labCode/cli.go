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
		fmt.Println(CliHandler(strings.Fields(read), kademlia, exit, reader))
	}
}

func CliHandler(s []string, kademlia *Kademlia, exit chan int, reader *bufio.Reader) string {
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
		if data.data != "" {
			return "get " + s[1] + " found the data: " + data.data
		}
		return "get " + s[1] + " did not find the data.\n 'get' found the contactsFound contacts: " + ContactsString(contacts)
	case "exit":
		if len(s) != 1 {
			return "Expected 0 arguments for command 'exit'"
		} else {
			exit <- 0
			return ""
		}
	case "refresh":
		if len(s) != 2 {
			return "Expected exactly 1 argument for command 'refresh'"
		} else if len(s[1]) != 40 {
			return "Argument provided to get is not a SHA-1 hash"
		}
		ch := make(chan int)
		kademlia.refreshMap[s[1]] = ch
		go kademlia.Refresh(s[1], ch)
		fmt.Println("Refreshing " + s[1])
		for {
			fmt.Print(">")
			read, _ := reader.ReadString('\n')
			fmt.Println(CliHandler(strings.Fields(read), kademlia, exit, reader))
		}
	case "forget":
		if len(s) != 2 {
			return "Expected exactly 1 argument for command 'forget'"
		} else if len(s[1]) != 40 {
			return "Argument provided to get is not a SHA-1 hash"
		}
		ch, ok := kademlia.refreshMap[s[1]]
		if ok {
			ch <- 0
			return "Forgot the data with hash: " + s[1]
		}
		return "Not refreshing any data with hash: " + s[1]
	default:
		return "Operation: " + operation + " not found"
	}
}
