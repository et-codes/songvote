package songvote

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/et-codes/songvote/internal/httplogger"
)

type SongStore interface {
	GetSong(id int64) (Song, error)
	GetSongs() Songs
	AddSong(song Song) (int64, error)
	DeleteSong(id int64) error
	UpdateSong(id int64, song Song) error
	AddVote(id int64) error
	Veto(id int64) error
}

type Server struct {
	store SongStore
	http.Handler
}

// NewServer returns a reference to an initialized Server.
func NewServer(store SongStore) *Server {
	s := new(Server)

	s.store = store

	router := http.NewServeMux()
	router.Handle("/songs/vote/", http.HandlerFunc(s.handleVote))   // POST
	router.Handle("/songs/veto/", http.HandlerFunc(s.handleVeto))   // POST
	router.Handle("/songs/", http.HandlerFunc(s.handleSongsWithID)) // GET|PATCH|DELETE
	router.Handle("/songs", http.HandlerFunc(s.handleSongs))        // GET|POST

	loggingRouter := httplogger.New(router)

	s.Handler = loggingRouter

	return s
}

// handleSongs routes requests to "/songs" depending on request type.
//
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
//
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
//
// Allowable methods:
//   - POST: vote for a song
func (s *Server) handleVote(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// TODO implement this method
		log.Println("POST vote not implemented")
		w.WriteHeader(http.StatusNotImplemented)
	default:
		code := http.StatusMethodNotAllowed
		message := fmt.Sprintf("Method %s not allowed", r.Method)
		writeError(w, code, message)
	}
}

// handleVote routes requests to "/songs/veto/{id}" depending on request type.
//
// Allowable methods:
//   - POST: veto a song
func (s *Server) handleVeto(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// TODO implement this method
		log.Println("POST veto not implemented")
		w.WriteHeader(http.StatusNotImplemented)
	default:
		code := http.StatusMethodNotAllowed
		message := fmt.Sprintf("Method %s not allowed", r.Method)
		writeError(w, code, message)
	}
}

func (s *Server) getAllSongs(w http.ResponseWriter, r *http.Request) {
	songs := s.store.GetSongs()

	out, err := MarshalJSON(songs)
	if err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Problem encoding songs to JSON: %v", err)
		writeError(w, code, message)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, out)
}

func (s *Server) getSong(w http.ResponseWriter, r *http.Request) {
	id, err := parseSongID(r.URL.Path, "/songs/")
	if err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Problem parsing song ID: %v", err)
		writeError(w, code, message)
		return
	}

	song, err := s.store.GetSong(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	json, err := MarshalJSON(song)
	if err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Problem marshaling song to JSON: %v", err)
		writeError(w, code, message)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, json)
}

func (s *Server) addSong(w http.ResponseWriter, r *http.Request) {
	songToAdd := Song{}
	if err := UnmarshalJSON(r.Body, &songToAdd); err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Problem unmarshaling song: %v", err)
		writeError(w, code, message)
		return
	}

	id, err := s.store.AddSong(songToAdd)
	if err != nil {
		code := http.StatusConflict
		message := fmt.Sprintf("Could not add song: %v\n", err)
		writeError(w, code, message)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, id)
}

func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {
	id, err := parseSongID(r.URL.Path, "/songs/")
	if err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Problem parsing song ID: %v", err)
		writeError(w, code, message)
		return
	}

	if err := s.store.DeleteSong(id); err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Could not delete song: %v\n", err)
		writeError(w, code, message)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) updateSong(w http.ResponseWriter, r *http.Request) {
	id, err := parseSongID(r.URL.Path, "/songs/")
	if err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Problem parsing song ID: %v", err)
		writeError(w, code, message)
		return
	}

	songToUpdate, err := s.store.GetSong(id)
	if err != nil {
		code := http.StatusNotFound
		message := fmt.Sprintf("Unable to retreive song %d: %v", id, err)
		writeError(w, code, message)
		return
	}

	updatedSong := Song{}
	if err := UnmarshalJSON(r.Body, &updatedSong); err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Problem parsing request body: %v", err)
		writeError(w, code, message)
		return
	}

	songToUpdate.Name = updatedSong.Name
	songToUpdate.Artist = updatedSong.Artist
	songToUpdate.LinkURL = updatedSong.LinkURL

	err = s.store.UpdateSong(id, songToUpdate)
	if err != nil {
		code := http.StatusInternalServerError
		message := fmt.Sprintf("Error updating song %d: %v", id, err)
		writeError(w, code, message)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseSongID(path, prefix string) (int64, error) {
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
