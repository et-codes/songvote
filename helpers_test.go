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

var testSong = songvote.Song{
	Name:    "Mirror In The Bathroom",
	Artist:  "Oingo Boingo",
	LinkURL: "https://youtu.be/SHWrmIzgB5A?si=R96_BWKxol3i7kQe",
	Votes:   10,
	Vetoed:  false,
}

// setupStuite sets up test conditions and returns function to defer until
// after the test runs.
func setupSuite(t *testing.T) (
	teardownSuite func(t *testing.T),
	server *songvote.Server,
) {
	// Suppress logging to os.Stdout during tests.
	log.SetOutput(io.Discard)
	store := songvote.NewSQLStore(":memory:")
	server = songvote.NewServer(store)
	teardownSuite = func(t *testing.T) {
		// Restore logging to os.Stdout after tests.
		log.SetOutput(os.Stdout)
	}
	return
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
