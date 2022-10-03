package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

func main() {
	r, _ := hex.DecodeString("48654a")
	a := string(r)
	fmt.Printf("%s\n", r)
	fmt.Printf("%v\n", a)
	s := hex.EncodeToString(r)
	fmt.Printf("%T\n", s)
	fmt.Println(s)

	xd := []byte{0}
	fmt.Println(string(xd))
	inp := "hahaha     "
	fmt.Println([]byte(inp))
	inp, _, _ = strings.Cut(inp, " ")
	fmt.Println(len(inp))
	q := sha1.Sum([]byte("hahaha"))
	q2 := hex.EncodeToString(q[:])
	fmt.Println(q2)
}
