package songvote

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSongs(t *testing.T) {
	t.Run("returns a song", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/songs/would", nil)
		response := httptest.NewRecorder()

		SongVoteServer(response, request)

		got := response.Body.String()
		want := "Would?"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
