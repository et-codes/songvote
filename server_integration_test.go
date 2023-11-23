package songvote_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et-codes/songvote"
)

func TestGetAndAddSongs(t *testing.T) {
	t.Skip("Skipping to work on other stuff")
	store := songvote.InMemorySongStore{}
	server := songvote.NewServer(&store)
	songName := "Enjoy the Silence"

	server.ServeHTTP(httptest.NewRecorder(), newPostSongRequest(songName))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetSongRequest(songName))
	assertStatus(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), songName)
}
