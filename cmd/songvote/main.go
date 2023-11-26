package main

import (
	"log"
	"net/http"

	"github.com/et-codes/songvote"
)

const (
	port   = ":5050"
	dbPath = "./db/songvote.db"
)

func main() {
	store := songvote.NewSQLStore(dbPath)
	server := songvote.NewServer(store)
	log.Printf("Listening on port %s ...\n", port)
	log.Fatal(http.ListenAndServe(port, server))
}
