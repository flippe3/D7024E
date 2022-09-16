package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	r, _ := hex.DecodeString("48656c6c6f20476f7068657221")
	fmt.Printf("%s\n", r)
	s := hex.EncodeToString(r)
	fmt.Printf("%T\n", s)
	fmt.Println(s)

}
