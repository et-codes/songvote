package main

import (
	"log"
)

const (
	port   = ":5050"
	dbFile = "db/songvote.db"
)

func main() {
	store, err := NewStore(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	server := NewServer(port, store)
	log.Fatal(server.ListenAndServe())
}
