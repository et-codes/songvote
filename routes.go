package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(s *Server) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", s.index).Methods(http.MethodGet)
	router.HandleFunc("/api/user", s.createUser).Methods(http.MethodPost)
	router.HandleFunc("/api/login", s.loginUser).Methods(http.MethodPost)
	router.HandleFunc("/api/logout", s.logoutUser).Methods(http.MethodGet)

	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.Use(logRequests)

	return router
}
