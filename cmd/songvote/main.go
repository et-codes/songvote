package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/et-codes/songvote"
)

const (
	port   = ":5050"
	dbPath = "./db/songvote.db"
)

func main() {
	store := songvote.NewSQLiteStore(dbPath)
	server := songvote.NewServer(store)
	slog.Info("Listening...", "port", port)
	log.Fatal(http.ListenAndServe(port, server))
}
