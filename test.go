package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
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

	t := time.Now()
	ttl := 2
	time.Sleep(time.Duration(3e9))
	if t.Add(time.Duration(ttl * 1e9)).After(time.Now()) {
		fmt.Println("not expired")
	} else {
		fmt.Println("expired")
	}
}
