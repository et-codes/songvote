package songvote_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/et-codes/songvote"
)

var (
	testSong = songvote.Song{
		Name:    "Mirror In The Bathroom",
		Artist:  "Oingo Boingo",
		LinkURL: "https://youtu.be/SHWrmIzgB5A?si=R96_BWKxol3i7kQe",
		Votes:   10,
		Vetoed:  false,
	}
	testUser = songvote.User{
		Active:   true,
		Name:     "John Doe",
		Password: "p@ssword",
		Vetoes:   1,
	}
)

// setupSuite sets up test conditions and returns function to defer until
// after the test runs. Uses SQLiteStore for the store.
func setupSuite(t *testing.T) (
	teardownSuite func(t *testing.T),
	store *songvote.SQLiteStore,
	server *songvote.Server,
) {
	// Suppress logging to os.Stdout during tests.
	log.SetOutput(io.Discard)
	store = songvote.NewSQLiteStore(":memory:")
	server = songvote.NewServer(store)
	teardownSuite = func(t *testing.T) {
		// Restore logging to os.Stdout after tests.
		log.SetOutput(os.Stdout)
	}
	return
}

// setupStubSuite sets up test conditions and returns function to defer until
// after the test runs. Uses StubStore for the store.
func setupStubSuite(t *testing.T) (
	teardownSuite func(t *testing.T),
	store *songvote.StubStore,
	server *songvote.Server,
) {
	// Suppress logging to os.Stdout during tests.
	log.SetOutput(io.Discard)
	store = songvote.NewStubStore()
	server = songvote.NewServer(store)
	teardownSuite = func(t *testing.T) {
		// Restore logging to os.Stdout after tests.
		log.SetOutput(os.Stdout)
	}
	return
}

// newAddUserRequest returns an HTTP request to add the given user.
func newAddUserRequest(user songvote.User) *http.Request {
	json, err := songvote.MarshalJSON(user)
	log.Printf("%s  %#v", json, user)
	if err != nil {
		log.Fatalf("problem marshalling User JSON, %v", err)
	}
	bodyReader := bytes.NewBuffer([]byte(json))
	request, _ := http.NewRequest(http.MethodPost, "/users", bodyReader)
	return request
}

// populateWithUser adds a user to a server for testing purposes.
func populateWithUser(server *songvote.Server, user songvote.User) {
	request := newAddUserRequest(user)
	server.ServeHTTP(httptest.NewRecorder(), request)
}

// populateWithSong adds a song to a server for testing purposes.
func populateWithSong(server *songvote.Server, song songvote.Song) {
	request := newAddSongRequest(song)
	server.ServeHTTP(httptest.NewRecorder(), request)
}

func newGetSongRequest(id int64) *http.Request {
	url := fmt.Sprintf("/songs/%d", id)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newGetSongsRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/songs", nil)
	return request
}

func newAddSongRequest(song songvote.Song) *http.Request {
	json, err := songvote.MarshalJSON(song)
	if err != nil {
		log.Fatalf("problem marshalling Song JSON, %v", err)
	}
	bodyReader := bytes.NewBuffer([]byte(json))
	request, _ := http.NewRequest(http.MethodPost, "/songs", bodyReader)
	return request
}

func newDeleteSongRequest(id int64) *http.Request {
	url := fmt.Sprintf("/songs/%d", id)
	request, _ := http.NewRequest(http.MethodDelete, url, nil)
	return request
}

func newUpdateSongRequest(id int64, song songvote.Song) *http.Request {
	json, err := songvote.MarshalJSON(song)
	if err != nil {
		log.Fatalf("problem marshalling Song JSON, %v", err)
	}
	url := fmt.Sprintf("/songs/%d", id)
	bodyReader := bytes.NewBuffer([]byte(json))
	request, _ := http.NewRequest(http.MethodPatch, url, bodyReader)
	return request
}

func newVoteRequest(id int64) *http.Request {
	url := fmt.Sprintf("/songs/vote/%d", id)
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	return request
}

func newVetoRequest(id int64) *http.Request {
	url := fmt.Sprintf("/songs/veto/%d", id)
	request, _ := http.NewRequest(http.MethodPost, url, nil)
	return request
}
