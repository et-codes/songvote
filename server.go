package main

import (
	"encoding/json"
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
	router.HandleFunc("/api/user", s.createUser).Methods(http.MethodPost)

	router.Use(logRequests)

	slog.Info("Server listening", "port", port)
	return http.ListenAndServe(port, router)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if err := s.tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		slog.Error(err.Error())
	}
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	userReq := NewUserRequest{
		Name:     r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	id, err := s.store.CreateUser(userReq)
	if err != nil {
		writeError(w, err.(ServerError))
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := NewUserResponse{id}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming HTTP request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
