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
	GetSong(id int) Song
	GetSongs() []Song
	AddSong(song Song)
}

type Server struct {
	store SongStore
	http.Handler
}

func NewServer(store SongStore) *Server {
	s := new(Server)

	s.store = store

	router := http.NewServeMux()
	router.Handle("/songs", http.HandlerFunc(s.getAllSongs)) // GET all songs
	router.Handle("/song/", http.HandlerFunc(s.getSong))     // GET song by ID
	router.Handle("/song", http.HandlerFunc(s.addSong))      // POST song

	s.Handler = router

	return s
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
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idString := strings.TrimPrefix(r.URL.Path, "/song/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Printf("problem parsing song ID from %s, %v", idString, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	song := s.store.GetSong(id)

	if song.Name == "" {
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	songToAdd := Song{}
	if err := UnmarshalSong(r.Body, &songToAdd); err != nil {
		log.Printf("could not read message body %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.store.AddSong(songToAdd)
	w.WriteHeader(http.StatusAccepted)
}
