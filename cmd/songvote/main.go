package main

import (
	"log"
	"net/http"

	"github.com/et-codes/songvote"
)

const (
	port = ":5050"
)

func main() {
	server := songvote.NewServer(songvote.NewInMemorySongStore())
	log.Fatal(http.ListenAndServe(port, server))
}
