package main

import (
	"testing"
)

func TestCliHandler(t *testing.T) {
	exit := make(chan int)

	if CliHandler([]string{}, &Kademlia{}, exit) != "" {
		t.Errorf("Received error for nil input")
	}

	if CliHandler([]string{"put"}, &Kademlia{}, exit) != "Expected exactly 2 arguments for command 'put'" {
		t.Errorf("Should print: Expected exactly 2 arguments for command 'put'")
	}

	if CliHandler([]string{"get"}, &Kademlia{}, exit) != "Expected exactly 1 argument for command 'get'" {
		t.Errorf("Should print: Expected exactly 1 argument for command 'get'")
	}

	if CliHandler([]string{"get", "asd"}, &Kademlia{}, exit) != "Argument provided to get is not a SHA-1 hash" {
		t.Errorf("Should print: Argument provided to get is not a SHA-1 hash")
	}

	if CliHandler([]string{"exit", "fgddfggfdgdf"}, &Kademlia{}, exit) != "Expected 0 arguments for command 'exit'" {
		t.Errorf("Should print: Expected 0 arguments for command 'exit'")
	}
	go CliHandler([]string{"exit"}, &Kademlia{}, exit)
	<-exit

	if CliHandler([]string{"asd", "fgddfggfdgdf"}, &Kademlia{}, exit) != "operation: "+"asd"+" not found" {
		t.Errorf("Should print: operation: " + "asd" + " not found")
	}
}
