package songvote_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

//
// Integration tests verify function of real SongStore with the API.
//

func TestGetAddDeleteSongs(t *testing.T) {
	store := songvote.NewSQLSongStore(":memory:")
	server := songvote.NewServer(store)
	songToAdd := songvote.Song{
		Name:    "Mirror In The Bathroom",
		Artist:  "Oingo Boingo",
		LinkURL: "http://test.com",
		Votes:   10,
		Vetoed:  false,
	}

	t.Run("add song and get ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newAddSongRequest(songToAdd)
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "1"

		assert.Equal(t, response.Code, http.StatusAccepted)
		assert.Equal(t, got, want)
	})

	t.Run("returns 409 with duplicate song", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newAddSongRequest(songToAdd)
		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusConflict)
	})

	t.Run("get song with ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetSongRequest(1)
		server.ServeHTTP(response, request)

		got := songvote.Song{}
		_ = songvote.Unmarshal(response.Body, &got)
		want := songToAdd
		want.ID = 1

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, got, want)
	})

	t.Run("delete a song", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newDeleteSongRequest(1)
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusNoContent)

		response = httptest.NewRecorder()
		request = newGetSongRequest(1)
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}
