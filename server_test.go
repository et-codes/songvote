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

func (s StubSongStore) GetSong(name string) string {
	return s.songs[name]
}

func TestGetSongs(t *testing.T) {
	store := StubSongStore{
		map[string]string{
			"would": "Would",
			"zero":  "Zero",
		},
	}
	server := songvote.NewServer(store)

	t.Run("returns the song Would", func(t *testing.T) {
		request := newGetSongRequest("Would")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "Would")
	})

	t.Run("returns the song Zero", func(t *testing.T) {
		request := newGetSongRequest("Zero")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "Zero")
	})
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func newGetSongRequest(name string) *http.Request {
	urlString := fmt.Sprintf("/songs/%s", name)
	url := url.PathEscape(urlString)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}
