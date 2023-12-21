package main

import (
	"log/slog"
	"net/http"
	"text/template"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
)

// Server contains configuration for the server.
type Server struct {
	port           string
	tmpl           *template.Template
	store          *Store
	sessionManager *scs.SessionManager
}

// NewServer creates and configures a new server.
func NewServer(port string, store *Store) *Server {
	var tmpl = template.Must(template.ParseGlob("templates/*"))
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(store.db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	return &Server{
		port:           port,
		tmpl:           tmpl,
		store:          store,
		sessionManager: sessionManager,
	}
}

// ListenAndServe starts the web server.
func (s *Server) ListenAndServe() error {
	router := mux.NewRouter()
	router.HandleFunc("/", s.index).Methods(http.MethodGet)
	router.HandleFunc("/api/user", s.createUser).Methods(http.MethodPost)

	router.Use(logRequests)

	slog.Info("Server listening", "port", port)
	return http.ListenAndServe(port, s.sessionManager.LoadAndSave(router))
}

// index executes the index page template.
func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		UserID   int
		Username string
	}{
		UserID:   s.sessionManager.GetInt(r.Context(), "user_id"),
		Username: s.sessionManager.GetString(r.Context(), "username"),
	}

	if err := s.tmpl.ExecuteTemplate(w, "index.gohtml", data); err != nil {
		slog.Error(err.Error())
	}
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
