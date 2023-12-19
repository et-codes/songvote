package main

import (
	"log"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

const (
	port = ":5050"
)
//nolint
var tpl = template.Must(template.ParseGlob("templates/*"))

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", index)

	slog.Info("Server listening...", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}