package main

import (
	"log"
	"net/http"

	"github.com/et-codes/songvote"
)

const (
	port   = ":5050"
	dbPath = "./db/songs.db"
)

func main() {
	store := songvote.NewInMemorySongStore()
	server := songvote.NewServer(store)
	log.Fatal(http.ListenAndServe(port, server))
}
