package main

import (
	"log/slog"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type Server struct {
	port  string
	tpl   *template.Template
	store *Store
}

func NewServer(port string, store *Store) *Server {
	var tpl = template.Must(template.ParseGlob("templates/*"))

	return &Server{
		port:  port,
		tpl:   tpl,
		store: store,
	}
}

func (s *Server) ListenAndServe() error {
	router := mux.NewRouter()

	router.HandleFunc("/", s.index).Methods(http.MethodGet)

	router.Use(logRequests)

	slog.Info("Server listening", "port", port)
	return http.ListenAndServe(port, router)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if err := s.tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		slog.Error(err.Error())
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming HTTP request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
