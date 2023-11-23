package songvote

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type SongStore interface {
	GetSong(name string) string
	AddSong(name string)
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
	song := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/songs/"))

	songName := s.store.GetSong(song)

	if songName == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, songName)
}

func (s *Server) addSong(w http.ResponseWriter, r *http.Request) {
	song, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("could not read message body %v", err)
	}
	log.Printf("adding song %s\n", song)

	s.store.AddSong(string(song))
	w.WriteHeader(http.StatusAccepted)
}
