package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// Sets up a reader that splits the input into 2 strings
// [0] = operation, [1] = value
func CliParser(kademlia *Kademlia, exit chan int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		read, err := reader.ReadString('\n')
		input := strings.Fields(read)
		if err != nil {
			fmt.Println(err)
		}
		CliHandler(input, kademlia, exit)
	}
}

func CliHandler(s []string, kademlia *Kademlia, exit chan int) {
	if len(s) == 0 {
		return
	}
	switch operation := s[0]; operation {
	case "put":
		if len(s) > 2 {
			fmt.Println("Too many arguments provided to operation 'put'")
		} else {
			contacts := kademlia.Store(s[1])
			r := sha1.Sum([]byte(s[1]))
			fmt.Println("put ", s[1], " sent store RPCs with hash: ", hex.EncodeToString(r[:]), ", to contacts: ", ContactsString(contacts))
		}
	case "get":
		if len(s) > 2 {
			fmt.Println("Too many arguments provided to operation 'get'")
		} else if len(s[1]) != 40 {
			fmt.Println("Argument provided to get is not a SHA-1 hash")
		} else {
			data, contacts := kademlia.LookupData(s[1])
			if data != "" {
				fmt.Println("get ", s[1], " found the data: ", data)
			} else {
				fmt.Println("get ", s[1], " did not find the data.\n 'get' found the contactsFound contacts: ", ContactsString(contacts))
			}
		}
	case "exit":
		if len(s) > 1 {
			fmt.Println("Too many arguments provided to operation 'exit'")
		} else {
			exit <- 0
		}
	default:
		fmt.Println("operation: ", operation, " not found")
	}
}
