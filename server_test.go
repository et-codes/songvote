package songvote_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/et-codes/songvote"
)

type StubSongStore struct {
	songs     map[string]string
	postCalls []string
}

func (s *StubSongStore) GetSong(name string) string {
	songName := s.songs[name]
	return songName
}

func (s *StubSongStore) AddSong(name string) {
	s.postCalls = append(s.postCalls, name)
}

func TestGetSongs(t *testing.T) {
	store := StubSongStore{
		map[string]string{
			"would": "Would",
			"zero":  "Zero",
		},
		nil,
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

func TestStoreSongs(t *testing.T) {
	store := StubSongStore{
		map[string]string{},
		[]string{},
	}
	server := songvote.NewServer(&store)

	t.Run("stores song when POST", func(t *testing.T) {
		newSong := "Creep"

		request := newPostSongRequest(newSong)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.postCalls) != 1 {
			t.Errorf("got %d calls to AddSong, want %d", len(store.postCalls), 1)
		}

		if store.postCalls[0] != newSong {
			t.Errorf("did not store correct song, got %q, want %q", store.postCalls[0], newSong)
		}
	})
}

// Helper methods

func newGetSongRequest(name string) *http.Request {
	urlString := fmt.Sprintf("/songs/%s", name)
	url := url.PathEscape(urlString)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newPostSongRequest(name string) *http.Request {
	bodyReader := bytes.NewReader([]byte(name))
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
