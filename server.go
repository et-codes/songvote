package main

import (
	"log/slog"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type Server struct {
	port string
	tpl *template.Template
}

func NewServer(port string) *Server {
	var tpl = template.Must(template.ParseGlob("templates/*"))

	return &Server{
		port: port,
		tpl: tpl,
	}
}

func (s *Server) ListenAndServe() error {
	router := mux.NewRouter()

	router.HandleFunc("/", s.index).Methods(http.MethodGet)

	slog.Info("Server listening...", "port", port)
	return http.ListenAndServe(port, router)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if err := s.tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		slog.Error(err.Error())
	}
}