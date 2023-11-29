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
		ID:      1,
		Name:    "Mirror In The Bathroom",
		Artist:  "Oingo Boingo",
		LinkURL: "https://youtu.be/SHWrmIzgB5A?si=R96_BWKxol3i7kQe",
		Votes:   10,
		Vetoed:  false,
	}
	testUser = songvote.User{
		ID:       1,
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
	url := "/users"
	body := createResponseBody(user)
	request, _ := http.NewRequest(http.MethodPost, url, body)
	return request
}

func newGetUsersRequest() *http.Request {
	url := "/users"
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newGetUserRequest(id int64) *http.Request {
	url := fmt.Sprintf("/users/%d", id)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newDeleteUserRequest(id int64) *http.Request {
	url := fmt.Sprintf("/users/%d", id)
	request, _ := http.NewRequest(http.MethodDelete, url, nil)
	return request
}

func newUpdateUserRequest(id int64, user songvote.User) *http.Request {
	url := fmt.Sprintf("/users/%d", id)
	body := createResponseBody(user)
	request, _ := http.NewRequest(http.MethodPut, url, body)
	return request
}

func newGetSongRequest(id int64) *http.Request {
	url := fmt.Sprintf("/songs/%d", id)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newGetSongsRequest() *http.Request {
	url := "/songs"
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newAddSongRequest(song songvote.Song) *http.Request {
	url := "/songs"
	body := createResponseBody(song)
	request, _ := http.NewRequest(http.MethodPost, url, body)
	return request
}

func newDeleteSongRequest(id int64) *http.Request {
	url := fmt.Sprintf("/songs/%d", id)
	request, _ := http.NewRequest(http.MethodDelete, url, nil)
	return request
}

func newUpdateSongRequest(id int64, song songvote.Song) *http.Request {
	url := fmt.Sprintf("/songs/%d", id)
	body := createResponseBody(song)
	request, _ := http.NewRequest(http.MethodPut, url, body)
	return request
}

func newVoteRequest(vote songvote.Vote) *http.Request {
	url := "/songs/vote"
	body := createResponseBody(vote)
	request, _ := http.NewRequest(http.MethodPost, url, body)
	return request
}

func newGetVotesRequest(id int64) *http.Request {
	url := fmt.Sprintf("/songs/vote/%d", id)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newVetoRequest(veto songvote.Veto) *http.Request {
	url := "/songs/veto"
	body := createResponseBody(veto)
	request, _ := http.NewRequest(http.MethodPost, url, body)
	return request
}

func createResponseBody(obj any) *bytes.Buffer {
	json, err := songvote.MarshalJSON(obj)
	if err != nil {
		log.Fatalf("problem marshalling JSON, %v", err)
	}
	return bytes.NewBuffer([]byte(json))
}
