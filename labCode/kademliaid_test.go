package main

import (
	"testing"
)

func TestKademliaidRandomNewKademliaID(t *testing.T) {
	randomIdOne := NewRandomKademliaID()
	randomIdTwo := NewRandomKademliaID()

	if len(randomIdOne) != IDLength {
		t.Errorf("RandomNewKademliaId Wrong Length: %s", randomIdOne)
	}
	if randomIdOne == randomIdTwo {
		t.Errorf("Randomization not working id1: %s == id2: %s", randomIdOne, randomIdTwo)
	}
}

// This does not increase Less to 100%
func TestKademliaidLess(t *testing.T) {
	kadIdBig := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	kadIdSmall := NewKademliaID("1111111100000000000000000000000000000000")

	if kadIdSmall.Less(kadIdBig) == false {
		t.Errorf("Less is very broken it says %s < %s", kadIdBig, kadIdSmall)
	}
	if kadIdBig.Less(kadIdBig) == true {
		t.Errorf("Less is very broken it says %s < %s", kadIdBig, kadIdBig)
	}
}
