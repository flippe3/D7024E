package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func testPort(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to our Kademlia lab!\n")
	fmt.Println("CURL Request")
}

func getObjects(w http.ResponseWriter, r *http.Request, kademlia *Kademlia) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) != 40 {
		fmt.Printf("ID is wrong length\n")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data, _ := kademlia.LookupData(id)
	if data.data != "" {
		fmt.Printf("Successfully found the obj\n")
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		resp := make(map[string]string)
		resp["object"] = data.data

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Failed to create JSON %s", err)
		}
		w.Write(jsonResp)
		return
	} else {
		fmt.Printf("Couldn't find the data. \n")
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func postObjects(w http.ResponseWriter, r *http.Request, kademlia *Kademlia) {
	obj, _ := ioutil.ReadAll(r.Body)

	sha := sha1.Sum([]byte(string(obj)))
	ttl := 600 // default for now
	kademlia.Store(string(obj), ttl)
	fmt.Printf("Storing: %s\n", string(obj))

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["location"] = "/objects/" + hex.EncodeToString(sha[:])
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Failed to create JSON %s", err)
	}
	w.Write(jsonResp)
}

func handleRequests(kademlia *Kademlia) {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", testPort)
	myRouter.HandleFunc("/objects", func(w http.ResponseWriter, r *http.Request) { postObjects(w, r, kademlia) }).Methods("POST")
	myRouter.HandleFunc("/objects/{id}", func(w http.ResponseWriter, r *http.Request) { getObjects(w, r, kademlia) }).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
