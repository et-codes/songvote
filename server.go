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
	router.Handle("/songs/", http.HandlerFunc(s.handleSongsWithID)) // GET|PATCH|DELETE
	router.Handle("/songs", http.HandlerFunc(s.handleSongs))        // GET|POST

	router.Handle("/users/", http.HandlerFunc(s.handleUsersWithID)) // GET|PATCH|DELETE
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
		code := http.StatusMethodNotAllowed
		message := fmt.Sprintf("Method %s not allowed", r.Method)
		writeError(w, code, message)
	}
}

// handleUsersWithID routes requests to "/users/{id}" depending on request type.
// Allowable methods:
//   - GET:    get the user
//   - PATCH:  update the user
//   - DELETE: delete the user
func (s *Server) handleUsersWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getUser(w, r)
	// case http.MethodPost:
	//     // TODO
	default:
		code := http.StatusMethodNotAllowed
		message := fmt.Sprintf("Method %s not allowed", r.Method)
		writeError(w, code, message)
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
		code := http.StatusMethodNotAllowed
		message := fmt.Sprintf("Method %s not allowed", r.Method)
		writeError(w, code, message)
	}
}

// handleSongsWithID routes requests to "/songs/{id}" depending on request type.
// Allowable methods:
//   - GET:    get the song
//   - PATCH:  update the song
//   - DELETE: delete the song
func (s *Server) handleSongsWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getSong(w, r)
	case http.MethodPatch:
		s.updateSong(w, r)
	case http.MethodDelete:
		s.deleteSong(w, r)
	default:
		code := http.StatusMethodNotAllowed
		message := fmt.Sprintf("Method %s not allowed", r.Method)
		writeError(w, code, message)
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
		code := http.StatusMethodNotAllowed
		message := fmt.Sprintf("Method %s not allowed", r.Method)
		writeError(w, code, message)
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
		code := http.StatusMethodNotAllowed
		message := fmt.Sprintf("Method %s not allowed", r.Method)
		writeError(w, code, message)
	}
}

func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	userToAdd := User{}
	if err := UnmarshalJSON[User](r.Body, &userToAdd); err != nil {
		writeUnmarshalError(w, err)
		return
	}

	id, err := s.store.AddUser(userToAdd)
	if err != nil {
		writeConflictError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, id)
}

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := s.store.GetUsers()

	out, err := MarshalJSON(users)
	if err != nil {
		writeMarshalError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, out)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/users/")
	if err != nil {
		writeIDParseError(w, err)
		return
	}

	user, err := s.store.GetUser(id)
	if err != nil {
		writeNotFoundError(w, err)
		return
	}

	json, err := MarshalJSON(user)
	if err != nil {
		writeMarshalError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, json)
}

func (s *Server) getAllSongs(w http.ResponseWriter, r *http.Request) {
	songs := s.store.GetSongs()

	out, err := MarshalJSON(songs)
	if err != nil {
		writeMarshalError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, out)
}

func (s *Server) getSong(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/")
	if err != nil {
		writeIDParseError(w, err)
		return
	}

	song, err := s.store.GetSong(id)
	if err != nil {
		writeNotFoundError(w, err)
		return
	}

	json, err := MarshalJSON(song)
	if err != nil {
		writeMarshalError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, json)
}

func (s *Server) addSong(w http.ResponseWriter, r *http.Request) {
	songToAdd := Song{}
	if err := UnmarshalJSON[Song](r.Body, &songToAdd); err != nil {
		writeUnmarshalError(w, err)
		return
	}

	id, err := s.store.AddSong(songToAdd)
	if err != nil {
		writeConflictError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, id)
}

func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/")
	if err != nil {
		writeIDParseError(w, err)
		return
	}

	if err := s.store.DeleteSong(id); err != nil {
		code := http.StatusNotFound
		message := fmt.Sprintf("Could not delete song: %v\n", err)
		writeError(w, code, message)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) updateSong(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/")
	if err != nil {
		writeIDParseError(w, err)
		return
	}

	songToUpdate, err := s.store.GetSong(id)
	if err != nil {
		writeNotFoundError(w, err)
		return
	}

	updatedSong := Song{}
	if err := UnmarshalJSON[Song](r.Body, &updatedSong); err != nil {
		writeUnmarshalError(w, err)
		return
	}

	songToUpdate.Name = updatedSong.Name
	songToUpdate.Artist = updatedSong.Artist
	songToUpdate.LinkURL = updatedSong.LinkURL

	err = s.store.UpdateSong(id, songToUpdate)
	if err != nil {
		writeUpdateError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) addVote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/vote/")
	if err != nil {
		writeIDParseError(w, err)
		return
	}

	if err := s.store.AddVote(id); err != nil {
		writeUpdateError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) veto(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/songs/veto/")
	if err != nil {
		writeIDParseError(w, err)
		return
	}

	if err := s.store.Veto(id); err != nil {
		writeUpdateError(w, err)
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

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, NewError(code, message).ToJSON())
}

func writeNotFoundError(w http.ResponseWriter, err error) {
	code := http.StatusNotFound
	message := fmt.Sprintf("Problem retreiving: %v", err)
	writeError(w, code, message)
}

func writeIDParseError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	message := fmt.Sprintf("Problem parsing ID: %v", err)
	writeError(w, code, message)
}

func writeUpdateError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	message := fmt.Sprintf("Error updating: %v", err)
	writeError(w, code, message)
}

func writeUnmarshalError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	message := fmt.Sprintf("Problem parsing request body: %v", err)
	writeError(w, code, message)
}

func writeMarshalError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	message := fmt.Sprintf("Problem marshaling to JSON: %v", err)
	writeError(w, code, message)
}

func writeConflictError(w http.ResponseWriter, err error) {
	code := http.StatusConflict
	message := fmt.Sprintf("Resource already exists: %v", err)
	writeError(w, code, message)
}
