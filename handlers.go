package main

import (
	"log/slog"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		slog.Error(err.Error())
	}
}