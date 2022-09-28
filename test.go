package main

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/exp/slices"
)

func main() {
	r, _ := hex.DecodeString("48654a")
	a := string(r)
	fmt.Printf("%s\n", r)
	fmt.Printf("%v\n", a)
	s := hex.EncodeToString(r)
	fmt.Printf("%T\n", s)
	fmt.Println(s)

	v := []int{1, 2}
	fmt.Println(slices.Contains(v, 3))

}
