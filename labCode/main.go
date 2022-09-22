package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "172.22.0.2:81")
	if err != nil {
		fmt.Print("Couldn't dial")
	} else {
		fmt.Print("Dial worked")
		conn.Close()
	}
	ln, err := net.Listen("tcp", ":81")
	if err != nil {
		fmt.Print("Couldn't listen")
	} else {
		fmt.Print("Listening on: ", ln.Addr().String())
		conn, err := ln.Accept()
		if err != nil {
			fmt.Print("Couldn't accept")
		} else {
			fmt.Print("Connection accepted")
			conn.Close()
		}
	}

	/*
		network := Join()
		Listen(network.myIP, 81)

		reader := bufio.NewReader(os.Stdin)
		for {
			switch text, _ := reader.ReadString('\n'); text {
			case "join\n":

			case "put\n":

			case "get\n":

			case "exit\n":

			}
		}
	*/
}
