package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

const dbPath = "./db/songs_test.db"

func TestSongStore(t *testing.T) {
	newSong := songvote.Song{
		ID:      1,
		Name:    "Creep",
		Artist:  "Radiohead",
		LinkURL: "http://test.com",
		Votes:   12,
		Vetoed:  true,
	}
	store := songvote.NewSQLSongStore(dbPath)

	t.Run("adding song returns valid ID", func(t *testing.T) {
		id, err := store.AddSong(newSong)
		defer store.DeleteSong(id)
		assert.NoError(t, err)
		if id <= 0 {
			t.Errorf("got bad id %d", id)
		}
	})

	t.Run("can add and retreive a song", func(t *testing.T) {
		id, err := store.AddSong(newSong)
		defer store.DeleteSong(id)
		assert.NoError(t, err)
		got, err := store.GetSong(id)
		assert.NoError(t, err)
		assert.Equal(t, got, newSong)
	})

	t.Run("can't add duplicate song", func(t *testing.T) {
		id, err := store.AddSong(newSong)
		defer store.DeleteSong(id)
		assert.NoError(t, err)
		_, err = store.AddSong(newSong)
		assert.Error(t, err)
	})

	t.Run("gets all songs in store", func(t *testing.T) {
		t.Skip("this test pending development")
	})

	t.Run("deletes a song", func(t *testing.T) {
		id, _ := store.AddSong(newSong)
		err := store.DeleteSong(id)
		assert.NoError(t, err)
	})
}
