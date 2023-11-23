package songvote

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SongStore interface {
	GetSong(id int) Song
	AddSong(song Song)
}

type Server struct {
	store SongStore
}

func NewServer(store SongStore) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getSong(w, r)
	case http.MethodPost:
		s.addSong(w, r)
	}
}

func (s *Server) getSong(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/songs/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Printf("problem parsing song ID from %s, %v", idString, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	song := s.store.GetSong(id)

	if song.Name == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	json, err := song.Marshal()
	if err != nil {
		log.Printf("problem marshalling song to JSON, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, json)
}

func (s *Server) addSong(w http.ResponseWriter, r *http.Request) {
	json, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read message body %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	songToAdd := Song{}
	if err := UnmarshalSong(string(json), &songToAdd); err != nil {
		log.Printf("could not read message body %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	s.store.AddSong(songToAdd)
	w.WriteHeader(http.StatusAccepted)
}
