package songvote_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/et-codes/songvote"
)

func TestGetSongs(t *testing.T) {
	t.Run("returns a song", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/songs/would", nil)
		response := httptest.NewRecorder()

		songvote.Server(response, request)

		got := response.Body.String()
		want := "Would?"

		assertEqual(t, got, want)
	})
}

func assertEqual(t testing.TB, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
