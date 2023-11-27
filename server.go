package songvote

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/et-codes/songvote/internal/httplogger"
)

type Store interface {
	// Song methods
	AddSong(song Song) (int64, error)
	GetSongs() Songs
	GetSong(id int64) (Song, error)
	DeleteSong(id int64) error
	UpdateSong(id int64, song Song) error
	AddVote(id int64) error
	Veto(id int64) error

	// User methods
	AddUser(user User) (int64, error)
	GetUsers() Users
	GetUser(id int64) (User, error)
	DeleteUser(id int64) error
	UpdateUser(id int64, user User) error
}

type Server struct {
	store Store
	http.Handler
}

// NewServer returns a reference to an initialized Server.
func NewServer(store Store) *Server {
	s := new(Server)

	s.store = store

	router := http.NewServeMux()

	router.Handle("/songs/vote/", http.HandlerFunc(s.handleVote))   // POST
	router.Handle("/songs/veto/", http.HandlerFunc(s.handleVeto))   // POST
	router.Handle("/songs/", http.HandlerFunc(s.handleSongsWithID)) // GET|PUT|DELETE
	router.Handle("/songs", http.HandlerFunc(s.handleSongs))        // GET|POST

	router.Handle("/users/", http.HandlerFunc(s.handleUsersWithID)) // GET|PUT|DELETE
	router.Handle("/users", http.HandlerFunc(s.handleUsers))        // POST

	loggingRouter := httplogger.New(router)

	s.Handler = loggingRouter

	return s
}

// handleUsers routes requests to "/users" depending on request type.
// Allowable methods:
//   - GET:  get all users
//   - POST: add user
func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getAllUsers(w, r)
	case http.MethodPost:
		s.addUser(w, r)
	default:
		writeError(w, ErrMethod)
	}
}

// handleUsersWithID routes requests to "/users/{id}" depending on request type.
// Allowable methods:
//   - GET:    get the user
//   - PUT:    update the user
//   - DELETE: delete the user
func (s *Server) handleUsersWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getUser(w, r)
	case http.MethodPut:
		s.updateUser(w, r)
	case http.MethodDelete:
		s.deleteUser(w, r)
	default:
		writeError(w, ErrMethod)
	}
}

// handleSongs routes requests to "/songs" depending on request type.
// Allowable methods:
//   - GET:  get all songs
//   - POST: add song
func (s *Server) handleSongs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getAllSongs(w, r)
	case http.MethodPost:
		s.addSong(w, r)
	default:
		writeError(w, ErrMethod)
	}
}

// handleSongsWithID routes requests to "/songs/{id}" depending on request type.
// Allowable methods:
//   - GET:    get the song
//   - PUT:    update the song
//   - DELETE: delete the song
func (s *Server) handleSongsWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getSong(w, r)
	case http.MethodPut:
		s.updateSong(w, r)
	case http.MethodDelete:
		s.deleteSong(w, r)
	default:
		writeError(w, ErrMethod)
	}
}

// handleVote routes requests to "/songs/vote/{id}" depending on request type.
// Allowable methods:
//   - POST: vote for a song
func (s *Server) handleVote(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.addVote(w, r)
	default:
		writeError(w, ErrMethod)
	}
}

// handleVote routes requests to "/songs/veto/{id}" depending on request type.
// Allowable methods:
//   - POST: veto a song
func (s *Server) handleVeto(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.veto(w, r)
	default:
		writeError(w, ErrMethod)
	}
}

func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	userToAdd := User{}
	if err := UnmarshalJSON(r.Body, &userToAdd); err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	id, err := s.store.AddUser(userToAdd)
	if err != nil {
		writeError(w, ErrConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, id)
}

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := s.store.GetUsers()

	out, err := MarshalJSON(users)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, out)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/users/")
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	user, err := s.store.GetUser(id)
	if err != nil {
		writeError(w, ErrNotFound)
		return
	}

	json, err := MarshalJSON(user)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, json)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/users/")
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	if err := s.store.DeleteUser(id); err != nil {
		writeError(w, ServerError{http.StatusNotFound, err})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// updateUser updates the User with new name and/or password.
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/users/")
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	userToUpdate, err := s.store.GetUser(id)
	if err != nil {
		writeError(w, ErrNotFound)
		return
	}

	updatedUser := User{}
	if err := UnmarshalJSON(r.Body, &updatedUser); err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	userToUpdate.Name = updatedUser.Name
	userToUpdate.Password = updatedUser.Password

	err = s.store.UpdateUser(id, userToUpdate)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getAllSongs(w http.ResponseWriter, r *http.Request) {
	songs := s.store.GetSongs()

	out, err := MarshalJSON(songs)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, out)
}

func (s *Server) getSong(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/")
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	song, err := s.store.GetSong(id)
	if err != nil {
		writeError(w, ErrNotFound)
		return
	}

	json, err := MarshalJSON(song)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, json)
}

func (s *Server) addSong(w http.ResponseWriter, r *http.Request) {
	songToAdd := Song{}
	if err := UnmarshalJSON(r.Body, &songToAdd); err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	id, err := s.store.AddSong(songToAdd)
	if err != nil {
		writeError(w, ErrConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, id)
}

func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/")
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	if err := s.store.DeleteSong(id); err != nil {
		writeError(w, ServerError{http.StatusNotFound, err})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) updateSong(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/")
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	songToUpdate, err := s.store.GetSong(id)
	if err != nil {
		writeError(w, ErrNotFound)
		return
	}

	updatedSong := Song{}
	if err := UnmarshalJSON(r.Body, &updatedSong); err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	songToUpdate.Name = updatedSong.Name
	songToUpdate.Artist = updatedSong.Artist
	songToUpdate.LinkURL = updatedSong.LinkURL

	err = s.store.UpdateSong(id, songToUpdate)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) addVote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/vote/")
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	if err := s.store.AddVote(id); err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) veto(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/veto/")
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	if err := s.store.Veto(id); err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(path, prefix string) (int64, error) {
	idString := strings.TrimPrefix(path, prefix)
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
