package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	r, _ := hex.DecodeString("48654a")
	fmt.Printf("%v\n", r)
	s := hex.EncodeToString(r)
	fmt.Printf("%T\n", s)
	fmt.Println(s)

	fmt.Println(string(72))
}
