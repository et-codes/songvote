// The purpose of these tests is to ensure that Server properly receives
// the HTTP requests, calls the appropriate methods, and returns the
// correct data types.

package songvote_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

// StubSongStore implements the SongStore interface, and keeps track of
// the calls against its methods.
type StubSongStore struct {
	nextID            int64                   // next song ID to be used
	getSongCalls      []int64                 // calls to to GetSong
	getSongsCallCount int                     // count of calls to GetSongs
	addSongCalls      songvote.Songs          // calls to AddSong
	deleteSongCalls   []int64                 // calls to DeleteSong
	updateSongCalls   map[int64]songvote.Song // calls to UpdateSong
	addVoteCalls      []int64                 // calls to AddVote
	vetoCalls         []int64                 // calls to Veto
}

func NewStubSongStore() *StubSongStore {
	return &StubSongStore{
		nextID:          1,
		updateSongCalls: make(map[int64]songvote.Song),
	}
}

func (s *StubSongStore) GetSong(id int64) (songvote.Song, error) {
	s.getSongCalls = append(s.getSongCalls, id)
	if id >= 10 {
		return songvote.Song{}, fmt.Errorf("song ID %d not found", id)
	}
	return songvote.Song{}, nil
}

func (s *StubSongStore) AddSong(song songvote.Song) (int64, error) {
	song.ID = s.nextID
	s.nextID++
	s.addSongCalls = append(s.addSongCalls, song)
	return song.ID, nil
}

func (s *StubSongStore) DeleteSong(id int64) error {
	s.deleteSongCalls = append(s.deleteSongCalls, id)
	if id >= 10 {
		return fmt.Errorf("song ID %d not found", id)
	}
	return nil
}

func (s *StubSongStore) GetSongs() songvote.Songs {
	s.getSongsCallCount++
	return songvote.Songs{}
}

func (s *StubSongStore) UpdateSong(id int64, song songvote.Song) error {
	s.updateSongCalls[id] = song
	if id >= 10 {
		return fmt.Errorf("song ID %d not found", id)
	}
	return nil
}

func (s *StubSongStore) AddVote(id int64) error {
	s.addVoteCalls = append(s.addVoteCalls, id)
	if id >= 10 {
		return fmt.Errorf("song ID %d not found", id)
	}
	return nil
}

func (s *StubSongStore) Veto(id int64) error {
	s.vetoCalls = append(s.vetoCalls, id)
	if id >= 10 {
		return fmt.Errorf("song ID %d not found", id)
	}
	return nil
}

func TestGetAllSongsFromServer(t *testing.T) {
	store := NewStubSongStore()
	server := songvote.NewServer(store)
	request := newGetSongsRequest()
	response := httptest.NewRecorder()

	t.Run("get all songs", func(t *testing.T) {
		server.ServeHTTP(response, request)
		songs := songvote.Songs{}
		_ = songvote.UnmarshalJSON(response.Body, &songs)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, songs, songvote.Songs{})
		assert.Equal(t, store.getSongsCallCount, 1)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/songs", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestGetSongsFromServer(t *testing.T) {
	store := NewStubSongStore()
	server := songvote.NewServer(store)

	t.Run("get single song", func(t *testing.T) {
		request := newGetSongRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, len(store.getSongCalls), 1)
	})

	t.Run("returns 404 if song ID not found", func(t *testing.T) {
		request := newGetSongRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/songs/0", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestAddSongsToServer(t *testing.T) {
	store := NewStubSongStore()
	server := songvote.NewServer(store)

	t.Run("stores song when POST", func(t *testing.T) {
		newSong := songvote.Song{
			Name:   "Creep",
			Artist: "Radiohead",
		}

		request := newAddSongRequest(newSong)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusAccepted)
		assert.Equal(t, len(store.addSongCalls), 1)
		var id int64
		_ = songvote.UnmarshalJSON(response.Body, &id)
		assert.Equal(t, id, int64(1))

		if !newSong.Equal(store.addSongCalls[0]) {
			t.Errorf("did not store correct song, got %v, want %v",
				store.addSongCalls[0], newSong)
		}
	})
}

func TestDeleteSongFromServer(t *testing.T) {
	store := NewStubSongStore()
	server := songvote.NewServer(store)

	t.Run("delete song response", func(t *testing.T) {
		request := newDeleteSongRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, store.deleteSongCalls[0], int64(1))
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		request := newDeleteSongRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusInternalServerError)
		assert.Equal(t, store.deleteSongCalls[1], int64(10))
	})
}

func TestUpdateSongOnServer(t *testing.T) {
	store := NewStubSongStore()
	server := songvote.NewServer(store)
	song := songvote.Song{
		Name:   "Creep",
		Artist: "Radiohead",
	}

	t.Run("update song", func(t *testing.T) {
		request := newUpdateSongRequest(1, song)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, len(store.updateSongCalls), 1)
		if !song.Equal(store.updateSongCalls[1]) {
			t.Errorf("got %v, want %v", store.updateSongCalls[1], song)
		}
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		request := newUpdateSongRequest(10, song)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}

func TestAddVoteOnServer(t *testing.T) {
	store := NewStubSongStore()
	server := songvote.NewServer(store)

	t.Run("update vote count", func(t *testing.T) {
		request := newVoteRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, store.addVoteCalls[0], int64(1))
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/songs/vote/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		request := newVoteRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusInternalServerError)
	})
}

func TestVetoSongOnServer(t *testing.T) {
	store := NewStubSongStore()
	server := songvote.NewServer(store)

	t.Run("veto a song", func(t *testing.T) {
		request := newVetoRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, len(store.vetoCalls), 1)
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		request := newVetoRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/songs/veto/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

// Helper methods

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
