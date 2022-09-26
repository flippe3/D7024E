package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	kademlia := Kademlia{}
	kademlia.Join()
	kademlia.network.Listen(kademlia.dataMap)

	reader := bufio.NewReader(os.Stdin)
	for {
		switch text, _ := reader.ReadString('\n'); text {
		case "join\n":

		case "put\n":

		case "get\n":

		case "exit\n":

		}
	}

}
