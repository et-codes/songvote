// The purpose of these tests is to ensure that Server properly receives
// the HTTP requests, calls the appropriate methods, and returns the
// correct data types.

package songvote_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

func TestGetAllSongsFromServer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	store := songvote.NewStubStore()
	server := songvote.NewServer(store)

	t.Run("get all songs", func(t *testing.T) {
		request := newGetSongsRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		songs := songvote.Songs{}
		_ = songvote.UnmarshalJSON(response.Body, &songs)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, songs, songvote.Songs{})
		assert.Equal(t, store.GetSongsCallCount, 1)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/songs", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestGetSongsFromServer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	store := songvote.NewStubStore()
	server := songvote.NewServer(store)

	t.Run("get single song", func(t *testing.T) {
		request := newGetSongRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, len(store.GetSongCalls), 1)
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
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	store := songvote.NewStubStore()
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
		assert.Equal(t, len(store.AddSongCalls), 1)
		var id int64
		_ = songvote.UnmarshalJSON(response.Body, &id)
		assert.Equal(t, id, int64(1))

		if !newSong.Equal(store.AddSongCalls[0]) {
			t.Errorf("did not store correct song, got %v, want %v",
				store.AddSongCalls[0], newSong)
		}
	})
}

func TestDeleteSongFromServer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	store := songvote.NewStubStore()
	server := songvote.NewServer(store)

	t.Run("delete song response", func(t *testing.T) {
		request := newDeleteSongRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, store.DeleteSongCalls[0], int64(1))
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		request := newDeleteSongRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusInternalServerError)
		assert.Equal(t, store.DeleteSongCalls[1], int64(10))
	})
}

func TestUpdateSongOnServer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	store := songvote.NewStubStore()
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
		assert.Equal(t, len(store.UpdateSongCalls), 1)
		if !song.Equal(store.UpdateSongCalls[1]) {
			t.Errorf("got %v, want %v", store.UpdateSongCalls[1], song)
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
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	store := songvote.NewStubStore()
	server := songvote.NewServer(store)

	t.Run("update vote count", func(t *testing.T) {
		request := newVoteRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, store.AddVoteCalls[0], int64(1))
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
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	store := songvote.NewStubStore()
	server := songvote.NewServer(store)

	t.Run("veto a song", func(t *testing.T) {
		request := newVetoRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, len(store.VetoCalls), 1)
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
