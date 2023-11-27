package songvote_test

import (
	"bytes"
	"encoding/json"
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
	userTestDataFile = "./testdata/users.json"
	songTestDataFile = "./testdata/songs.json"

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

// populateWithUser adds a user to a server for testing purposes.
func populateWithUser(server *songvote.Server, user songvote.User) {
	request := newAddUserRequest(user)
	server.ServeHTTP(httptest.NewRecorder(), request)
}

// populateWithUsers adds users to the server from JSON file.
func populateWithUsers(server *songvote.Server, path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not open file %s: %v", path, err)
	}
	users := songvote.Users{}
	if err = json.Unmarshal(file, &users); err != nil {
		log.Fatalf("Could not unmarshal JSON from %s: %v", path, err)
	}
	for _, user := range users {
		populateWithUser(server, user)
	}
}

// populateWithSong adds a song to a server for testing purposes.
func populateWithSong(server *songvote.Server, song songvote.Song) {
	request := newAddSongRequest(song)
	server.ServeHTTP(httptest.NewRecorder(), request)
}

// populateWithSongs adds songs to the server from JSON file.
func populateWithSongs(server *songvote.Server, path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not open file %s: %v", path, err)
	}
	songs := songvote.Songs{}
	if err = json.Unmarshal(file, &songs); err != nil {
		log.Fatalf("Could not unmarshal JSON from %s: %v", path, err)
	}
	for _, song := range songs {
		populateWithSong(server, song)
	}
}

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

func newGetUsersRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/users", nil)
	return request
}

func newGetUserRequest(id int64) *http.Request {
	url := fmt.Sprintf("/users/%d", id)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
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
