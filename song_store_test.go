package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

const dbPath = "./db/songs_test.db"

var store = songvote.NewSQLSongStore(dbPath)

func TestSongStore(t *testing.T) {
	newSong := songvote.Song{
		ID:      1,
		Name:    "Fake Song #27",
		Artist:  "Fake Artist",
		LinkURL: "http://test.com",
		Votes:   12,
		Vetoed:  true,
	}

	t.Run("adding song returns valid ID", func(t *testing.T) {
		id, err := addSong(t, newSong)
		assert.NoError(t, err)
		if id <= 0 {
			t.Errorf("got bad id %d", id)
		}
	})

	t.Run("can add and retreive a song", func(t *testing.T) {
		id, err := addSong(t, newSong)
		assert.NoError(t, err)
		got, err := store.GetSong(id)
		assert.NoError(t, err)
		assert.Equal(t, got, newSong)
	})

	t.Run("can't add duplicate song", func(t *testing.T) {
		_, err := addSong(t, newSong)
		assert.NoError(t, err)
		_, err = addSong(t, newSong)
		assert.Error(t, err)
	})

	t.Run("gets all songs in store", func(t *testing.T) {
		t.Skip("this test pending development")
	})

	t.Run("deletes a song", func(t *testing.T) {
		id, _ := addSong(t, newSong)
		err := store.DeleteSong(id)
		assert.NoError(t, err)
	})
}

// addSong is a helper function that calls the store's AddSong method and
// deletes the song after each test is concluded.
func addSong(t testing.TB, song songvote.Song) (int64, error) {
	id, err := store.AddSong(song)
	t.Cleanup(func() {
		_ = store.DeleteSong(id)
	})
	return id, err
}
