package main

import "log"

const (
	port = ":5050"
)

func main() {
	server := NewServer(port)
	log.Fatal(server.ListenAndServe())
}