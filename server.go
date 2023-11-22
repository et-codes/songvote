package songvote

import (
	"fmt"
	"net/http"
	"strings"
)

type SongStore interface {
	GetSong(name string) string
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
	song := strings.TrimPrefix(r.URL.Path, "/songs/")
	song = strings.ToLower(song)
	fmt.Fprint(w, s.store.GetSong(song))
}
