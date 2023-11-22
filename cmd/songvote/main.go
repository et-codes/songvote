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
	handler := http.HandlerFunc(songvote.Server)
	log.Fatal(http.ListenAndServe(port, handler))
}
