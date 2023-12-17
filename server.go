package songvote

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Store interface {
	// Song methods
	AddSong(song Song) (int64, error)
	GetSongs() Songs
	GetSong(id int64) (Song, error)
	DeleteSong(id int64) error
	UpdateSong(id int64, song Song) error
	AddVote(vote Vote) error
	GetVotesForSong(id int64) (Votes, error)
	Veto(veto Veto) error

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

var tpl *template.Template

// NewServer returns a reference to an initialized Server.
func NewServer(store Store) *Server {
	tpl = template.Must(template.ParseGlob("templates/*"))

	s := new(Server)
	s.store = store

	r := mux.NewRouter()

	r.HandleFunc("/songs/vote/{id}", s.getVotes).Methods(http.MethodGet)
	r.HandleFunc("/songs/vote", s.addVote).Methods(http.MethodPost)
	r.HandleFunc("/songs/veto", s.veto).Methods(http.MethodPost)
	r.HandleFunc("/songs/{id}", s.getSong).Methods(http.MethodGet)
	r.HandleFunc("/songs/{id}", s.updateSong).Methods(http.MethodPut)
	r.HandleFunc("/songs/{id}", s.deleteSong).Methods(http.MethodDelete)
	r.HandleFunc("/songs", s.getAllSongs).Methods(http.MethodGet)
	r.HandleFunc("/songs", s.addSong).Methods(http.MethodPost)

	r.HandleFunc("/users/{id}", s.getUser).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", s.updateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", s.deleteUser).Methods(http.MethodDelete)
	r.HandleFunc("/users", s.getAllUsers).Methods(http.MethodGet)
	r.HandleFunc("/users", s.addUser).Methods(http.MethodPost)

	r.HandleFunc("/", s.index).Methods(http.MethodGet)

	r.Use(logRequests)

	s.Handler = r

	return s
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	songs := s.store.GetSongs()

	if err := tpl.ExecuteTemplate(w, "index.gohtml", songs); err != nil {
		slog.Error("error executing template", err)
	}
}

func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	userToAdd := User{}
	if err := UnmarshalJSON(r.Body, &userToAdd); err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	id, err := s.store.AddUser(userToAdd)
	if err != nil {
		log.Printf("Error adding user: %v", err)
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
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, out)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
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
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, json)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	if err := s.store.DeleteUser(id); err != nil {
		writeError(w, ServerError{http.StatusNotFound, err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// updateUser updates the User with new name and/or password.
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
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
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	userToUpdate.Name = updatedUser.Name
	userToUpdate.Password = updatedUser.Password

	err = s.store.UpdateUser(id, userToUpdate)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getAllSongs(w http.ResponseWriter, r *http.Request) {
	songs := s.store.GetSongs()

	out, err := MarshalJSON(songs)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, out)
}

func (s *Server) getSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
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
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, json)
}

func (s *Server) addSong(w http.ResponseWriter, r *http.Request) {
	songToAdd := Song{
		Name: r.FormValue("songName"),
		Artist: r.FormValue("songArtist"),
		LinkURL: r.FormValue("songLinkURL"),
	}

	_, err := s.store.AddSong(songToAdd)
	if err != nil {
		writeError(w, ErrConflict)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	if err := s.store.DeleteSong(id); err != nil {
		writeError(w, ServerError{http.StatusNotFound, err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) updateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
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
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	songToUpdate.Name = updatedSong.Name
	songToUpdate.Artist = updatedSong.Artist
	songToUpdate.LinkURL = updatedSong.LinkURL

	err = s.store.UpdateSong(id, songToUpdate)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) addVote(w http.ResponseWriter, r *http.Request) {
	vote := Vote{}
	err := UnmarshalJSON(r.Body, &vote)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	if err := s.store.AddVote(vote); err != nil {
		switch err.Error() {
		case "user already voted for this song", "user is inactive and cannot perform this action":
			writeError(w, ErrConflict)
		default:
			writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getVotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeError(w, ErrIDParse)
		return
	}

	votes, err := s.store.GetVotesForSong(id)
	if err != nil {
		writeError(w, ErrNotFound)
		return
	}

	out, err := MarshalJSON(votes)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, out)
}

func (s *Server) veto(w http.ResponseWriter, r *http.Request) {
	veto := Veto{}
	err := UnmarshalJSON(r.Body, &veto)
	if err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	if err := s.store.Veto(veto); err != nil {
		writeError(w, ServerError{http.StatusInternalServerError, err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body string

		// Check if the body contains anything.
		if r.ContentLength > 0 {
			// Read body contents.
			buf, _ := io.ReadAll(r.Body)
			body = string(buf)

			// Put body contents into a reader and add it back to the request.
			reader := io.NopCloser(bytes.NewBuffer(buf))
			r.Body = reader
		}

		// Write log message.
		log.Printf("%s - %s (%s) - %s", r.Method, r.URL.Path, r.RemoteAddr, body)

		next.ServeHTTP(w, r)
	})
}
