package songvote_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

// These integration tests verify function of real Store through the API.

func TestAddSongs(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	t.Run("add song and get ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newAddSongRequest(testSong)

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "1"

		assert.Equal(t, response.Code, http.StatusAccepted)
		assert.Equal(t, got, want)
	})

	t.Run("returns 409 with duplicate song", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newAddSongRequest(testSong)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusConflict)
	})
}

func TestGetSong(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("get song with ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetSongRequest(1)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)

		got := songvote.Song{}
		err := songvote.UnmarshalJSON[songvote.Song](response.Body, &got)
		assert.NoError(t, err)

		want := testSong
		want.ID = 1

		assert.Equal(t, got, want)
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetSongRequest(999)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("get all songs", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetSongsRequest()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)

		songs := songvote.Songs{}
		err := songvote.UnmarshalJSON[songvote.Songs](response.Body, &songs)
		assert.NoError(t, err)

		assert.Equal(t, len(songs), 1)
		if !songs[0].Equal(testSong) {
			t.Errorf("want %v, got %v", testSong, songs[0])
		}
	})
}

func TestDeleteSong(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("delete a song and cannot retreive it", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newDeleteSongRequest(1)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)

		response = httptest.NewRecorder()
		request = newGetSongRequest(1)
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns error with unknown ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newDeleteSongRequest(999)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}
