package main

import (
	"log/slog"
	"net/http"
	"text/template"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Server contains configuration for the server.
type Server struct {
	port           string              // port number
	tmpl           *template.Template  // parsed templates
	store          *Store              // data storage
	sessionManager *scs.SessionManager // session manager
}

// PageData holds data being passed to a page template.
type PageData struct {
	UserID   int
	Username string
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
	router.HandleFunc("/", s.render("index.gohtml")).Methods(http.MethodGet)
	router.HandleFunc("/login", s.render("login.gohtml")).Methods(http.MethodGet)
	router.HandleFunc("/register", s.render("register.gohtml")).Methods(http.MethodGet)
	router.HandleFunc("/api/user", s.createUser).Methods(http.MethodPost)
	router.HandleFunc("/api/login", s.loginUser).Methods(http.MethodPost)
	router.HandleFunc("/api/logout", s.logoutUser).Methods(http.MethodGet)

	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.Use(logRequests)

	slog.Info("Server listening", "port", port)
	return http.ListenAndServe(port, s.sessionManager.LoadAndSave(router))
}

// render returns a handler function that executes the given template.
func (s *Server) render(template string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			UserID:   s.sessionManager.GetInt(r.Context(), "user_id"),
			Username: s.sessionManager.GetString(r.Context(), "username"),
		}
		if err := s.tmpl.ExecuteTemplate(w, template, data); err != nil {
			slog.Error(err.Error())
		}
	}
}

// logoutUser logs out the user by clearing session data.
func (s *Server) logoutUser(w http.ResponseWriter, r *http.Request) {
	username := s.sessionManager.Get(r.Context(), "username")
	id := s.sessionManager.Get(r.Context(), "user_id")

	if err := s.sessionManager.Clear(r.Context()); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("Logged out user", "user", username, "ID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
