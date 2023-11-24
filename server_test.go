package songvote_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

type StubSongStore struct {
	songs     []songvote.Song
	nextId    int
	postCalls []songvote.Song
}

func (s *StubSongStore) GetSong(id int) (songvote.Song, error) {
	for _, s := range s.songs {
		if s.ID == id {
			return s, nil
		}
	}
	return songvote.Song{}, fmt.Errorf("song ID %d not found", id)
}

func (s *StubSongStore) AddSong(song songvote.Song) (int, error) {
	song.ID = s.nextId
	s.nextId++
	s.postCalls = append(s.postCalls, song)
	return song.ID, nil
}

func (s *StubSongStore) GetSongs() []songvote.Song {
	return s.songs
}

func TestGetAllSongs(t *testing.T) {
	store := newPopulatedSongStore()
	server := songvote.NewServer(store)
	request := newGetSongsRequest()
	response := httptest.NewRecorder()

	t.Run("can get all songs", func(t *testing.T) {
		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assertSongListsEqual(t, response.Body, store.songs)
	})

	t.Run("returns empty array if empty store", func(t *testing.T) {
		store = newEmptySongStore()
		server = songvote.NewServer(store)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assertSongListsEqual(t, response.Body, store.songs)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ = http.NewRequest(http.MethodPost, "/songs", nil)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestGetSongs(t *testing.T) {
	store := newPopulatedSongStore()
	server := songvote.NewServer(store)

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/song/0", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)

		request, _ = http.NewRequest(http.MethodGet, "/song", nil)
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})

	t.Run("returns the song Would?", func(t *testing.T) {
		request := newGetSongRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		want, err := songvote.MarshalSong(store.songs[0])
		assert.NoError(t, err)

		got := response.Body.String()
		assert.Equal(t, got, want)
	})

	t.Run("returns the song Zero", func(t *testing.T) {
		request := newGetSongRequest(2)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		want, err := songvote.MarshalSong(store.songs[1])
		assert.NoError(t, err)

		got := response.Body.String()
		assert.Equal(t, got, want)
	})

	t.Run("returns 404 if song ID not found", func(t *testing.T) {
		request := newGetSongRequest(3)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreSongs(t *testing.T) {
	store := newEmptySongStore()
	server := songvote.NewServer(store)

	t.Run("stores song when POST", func(t *testing.T) {
		newSong := songvote.Song{
			Name:   "Creep",
			Artist: "Radiohead",
		}

		request := newPostSongRequest(newSong)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusAccepted)

		if len(store.postCalls) != 1 {
			t.Errorf("got %d calls to AddSong, want %d", len(store.postCalls), 1)
		}

		if !newSong.Equal(store.postCalls[0]) {
			t.Errorf("did not store correct song, got %v, want %v", store.postCalls[0], newSong)
		}
	})
}

// Helper methods

func newGetSongRequest(id int) *http.Request {
	url := fmt.Sprintf("/song/%d", id)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newGetSongsRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/songs", nil)
	return request
}

func newPostSongRequest(song songvote.Song) *http.Request {
	json, err := songvote.MarshalSong(song)
	if err != nil {
		log.Fatalf("problem marshalling Song JSON, %v", err)
	}
	bodyReader := bytes.NewReader([]byte(json))
	request, _ := http.NewRequest(http.MethodPost, "/song", bodyReader)
	return request
}

func newPopulatedSongStore() *StubSongStore {
	return &StubSongStore{
		songs: []songvote.Song{
			{ID: 1, Name: "Would?", Artist: "Alice in Chains"},
			{ID: 2, Name: "Zero", Artist: "The Smashing Pumpkins"},
		},
		postCalls: []songvote.Song{},
	}
}

func newEmptySongStore() *StubSongStore {
	return &StubSongStore{
		songs:     []songvote.Song{},
		postCalls: []songvote.Song{},
	}
}

// Custom assertions

func assertSongListsEqual(t testing.TB, body *bytes.Buffer, songs []songvote.Song) {
	t.Helper()
	want := []songvote.Song{}
	err := json.NewDecoder(body).Decode(&want)
	if err != nil {
		t.Errorf("could not decode response to JSON %v", err)
	}
	assert.Equal(t, want, songs)
}
