package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Server contains configuration for the server.
type Server struct {
	port           string              // port number
	store          *Store              // data storage
	sessionManager *scs.SessionManager // session manager
}

// NewServer creates and configures a new server.
func NewServer(port string, store *Store) *Server {
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(store.db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	return &Server{
		port:           port,
		store:          store,
		sessionManager: sessionManager,
	}
}

// ListenAndServe starts the web server.
func (s *Server) ListenAndServe() error {
	router := mux.NewRouter()

	// Template routes
	router.Handle("/", templ.Handler(index())).Methods(http.MethodGet)

	// API routes
	router.HandleFunc("/api/user", s.createUser).Methods(http.MethodPost)
	router.HandleFunc("/api/user", s.getUsers).Methods(http.MethodGet)
	router.HandleFunc("/api/user/{id}", s.getUser).Methods(http.MethodGet)
	router.HandleFunc("/api/user/{id}", s.deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/api/user/{id}", s.updateUser).Methods(http.MethodPut)
	router.HandleFunc("/api/login", s.loginUser).Methods(http.MethodPost)
	router.HandleFunc("/api/logout", s.logoutUser).Methods(http.MethodGet)

	// Middleware
	router.Use(logRequests)

	// Start the server
	slog.Info("Server listening", "port", port)
	return http.ListenAndServe(port, s.sessionManager.LoadAndSave(router))
}

// getUsers returns a list of all users with their ids.
func (s *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.GetUsers()
	if err != nil {
		writeError(w, NewServerError(http.StatusInternalServerError, err.Error()))
	}

	writeJSON(w, http.StatusOK, users)
}

// getUser returns the user with the given id.
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, NewServerError(http.StatusInternalServerError, err.Error()))
	}

	user, err := s.store.GetUserByID(userID)
	if err != nil {
		writeError(w, ErrNotFound)
	}

	writeJSON(w, http.StatusOK, user)
}

// updateUser updates a user.
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, NewServerError(http.StatusInternalServerError, err.Error()))
	}

	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		writeError(w, NewServerError(http.StatusInternalServerError, err.Error()))
	}
	user.ID = id

	if err := s.store.UpdateUser(user); err != nil {
		serverError := err.(ServerError)
		writeError(w, serverError)
	}

	writeJSON(w, http.StatusOK, user)
}

// deleteUser deletes the user with the given ID.
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeError(w, NewServerError(http.StatusInternalServerError, err.Error()))
	}

	if err := s.store.DeleteUser(id); err != nil {
		slog.Error("error deleting user", "id", id, "error", err.Error())
		writeError(w, ErrNotFound)
	}

	writeJSON(w, http.StatusNoContent, nil)
}

// logoutUser logs out the user by clearing session data.
func (s *Server) logoutUser(w http.ResponseWriter, r *http.Request) {
	username := s.sessionManager.Get(r.Context(), "username")
	id := s.sessionManager.Get(r.Context(), "user_id")

	if err := s.sessionManager.Clear(r.Context()); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("Logged out user", "user", username, "ID", id)

	writeJSON(w, http.StatusNoContent, nil)
}

// loginUser processes requests to log in an existing user.
func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := s.store.GetUserByName(username)
	if err != nil {
		writeError(w, ErrNotFound)
		return
	}

	if user.Inactive {
		writeError(w, NewServerError(http.StatusUnauthorized, "user is inactive"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		writeError(w, NewServerError(http.StatusUnauthorized,
			"incorrect username and/or password"))
		return
	}

	s.sessionManager.Put(r.Context(), "user_id", user.ID)
	s.sessionManager.Put(r.Context(), "username", user.Name)
	slog.Info("Logged in user", "user", user.Name, "ID", user.ID)

	writeJSON(w, http.StatusNoContent, nil)
}

// createUser processes requests to create a new user.
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

	s.sessionManager.Put(r.Context(), "user_id", id)
	s.sessionManager.Put(r.Context(), "username", userReq.Name)

	newUser := NewUserResponse{id, userReq.Name}

	writeJSON(w, http.StatusCreated, newUser)
}

// writeJSON encodes v into a JSON object and writes it to the response writer
// with the provided status code in the header.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			writeError(w, err.(ServerError))
		}
	}
}
