package main

import "log"

const (
	port = ":5050"
)

func main() {
	store, err := NewStore("db/songvote.db")
	if err != nil {
		log.Fatal(err)
	}
	server := NewServer(port, store)
	log.Fatal(server.ListenAndServe())
}
