package songvote_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/et-codes/songvote"
)

type StubSongStore struct {
	songs     []songvote.Song
	postCalls []songvote.Song
}

func (s *StubSongStore) GetSong(id int) songvote.Song {
	for _, s := range s.songs {
		if s.ID == id {
			return s
		}
	}
	return songvote.Song{}
}

func (s *StubSongStore) AddSong(song songvote.Song) {
	s.postCalls = append(s.postCalls, song)
}

func TestGetSongs(t *testing.T) {
	store := StubSongStore{
		songs: []songvote.Song{
			{ID: 0, Name: "Would?", Artist: "Alice in Chains"},
			{ID: 1, Name: "Zero", Artist: "The Smashing Pumpkins"},
		},
	}
	server := songvote.NewServer(&store)

	t.Run("returns the song Would?", func(t *testing.T) {
		request := newGetSongRequest(0)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)

		want, err := store.songs[0].Marshal()
		assertNoError(t, err)

		got := response.Body.String()
		assertResponseBody(t, got, want)
	})

	t.Run("returns the song Zero", func(t *testing.T) {
		request := newGetSongRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)

		want, err := store.songs[1].Marshal()
		assertNoError(t, err)

		got := response.Body.String()
		assertResponseBody(t, got, want)
	})

	t.Run("returns 404 if song not found", func(t *testing.T) {
		request := newGetSongRequest(3)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreSongs(t *testing.T) {
	store := StubSongStore{
		songs:     []songvote.Song{},
		postCalls: []songvote.Song{},
	}
	server := songvote.NewServer(&store)

	t.Run("stores song when POST", func(t *testing.T) {
		newSong := songvote.Song{
			Name:   "Creep",
			Artist: "Radiohead",
		}

		request := newPostSongRequest(newSong)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.postCalls) != 1 {
			t.Errorf("got %d calls to AddSong, want %d", len(store.postCalls), 1)
		}

		if store.postCalls[0] != newSong {
			t.Errorf("did not store correct song, got %v, want %v", store.postCalls[0], newSong)
		}
	})
}

func TestGetSongByName(t *testing.T) {

}

// Helper methods

func newGetSongRequest(id int) *http.Request {
	url := fmt.Sprintf("/songs/%d", id)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newPostSongRequest(song songvote.Song) *http.Request {
	json, err := song.Marshal()
	if err != nil {
		log.Fatalf("problem marshalling Song JSON, %v", err)
	}
	bodyReader := bytes.NewReader([]byte(json))
	request, _ := http.NewRequest(http.MethodPost, "/songs", bodyReader)
	return request
}

// Assertions

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("wrong message body, got %q, want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("wrong status code, got %d, want %d", got, want)
	}
}

func assertTrue(t testing.TB, got bool) {
	t.Helper()
	if !got {
		t.Errorf("got %t, wanted %t", got, true)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("got error but didn't want one, %v", err)
	}
}

func assertEqual(t testing.TB, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
