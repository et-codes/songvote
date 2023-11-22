package songvote_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/et-codes/songvote"
)

type StubSongStore struct {
	songs map[string]string
}

func (s *StubSongStore) GetSong(name string) string {
	return s.songs[name]
}

func TestGetSongs(t *testing.T) {
	store := StubSongStore{
		map[string]string{
			"would": "Would",
			"zero":  "Zero",
		},
	}
	server := songvote.NewServer(&store)

	t.Run("returns the song Would", func(t *testing.T) {
		request := newGetSongRequest("Would")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "Would")
	})

	t.Run("returns the song Zero", func(t *testing.T) {
		request := newGetSongRequest("Zero")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "Zero")
	})

	t.Run("returns 404 if song not found", func(t *testing.T) {
		request := newGetSongRequest("Jeepers Creepers")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

// Helper methods

func newGetSongRequest(name string) *http.Request {
	urlString := fmt.Sprintf("/songs/%s", name)
	url := url.PathEscape(urlString)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
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
