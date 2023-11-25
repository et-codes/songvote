package songvote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SongStore interface {
	GetSong(id int64) (Song, error)
	GetSongs() []Song
	AddSong(song Song) (int64, error)
	DeleteSong(id int64) error
}

type Server struct {
	store SongStore
	http.Handler
}

func NewServer(store SongStore) *Server {
	s := new(Server)

	s.store = store

	router := http.NewServeMux()
	router.Handle("/songs/vote/", http.HandlerFunc(s.handleVote))   // POST
	router.Handle("/songs/veto/", http.HandlerFunc(s.handleVeto))   // POST
	router.Handle("/songs/", http.HandlerFunc(s.handleSongsWithID)) // GET|PATCH|DELETE
	router.Handle("/songs", http.HandlerFunc(s.handleSongs))        // GET|POST

	loggingRouter := NewLoggerMiddleware(router)

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
		w.WriteHeader(http.StatusMethodNotAllowed)
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
		// TODO implement this method
		log.Println("PATCH song not implemented")
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodDelete:
		s.deleteSong(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
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
		w.WriteHeader(http.StatusMethodNotAllowed)
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
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) getAllSongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	songs := s.store.GetSongs()

	out := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(out).Encode(songs)
	if err != nil {
		log.Fatalf("problem encoding songs to JSON, %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, out)
}

func (s *Server) getSong(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/songs/")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Printf("problem parsing song ID from %s, %v", idString, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	song, err := s.store.GetSong(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	json, err := MarshalSong(song)
	if err != nil {
		log.Printf("problem marshalling song to JSON, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, json)
}

func (s *Server) addSong(w http.ResponseWriter, r *http.Request) {
	songToAdd := Song{}
	if err := UnmarshalSong(r.Body, &songToAdd); err != nil {
		log.Printf("could not read message body %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := s.store.AddSong(songToAdd)
	if err != nil {
		log.Printf("could not add song, %v\n", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, id)
}

func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/songs/")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Printf("problem parsing song ID from %s, %v", idString, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.store.DeleteSong(id); err != nil {
		log.Printf("could not delete song: %v\n", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
