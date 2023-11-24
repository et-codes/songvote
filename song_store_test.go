package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

func TestSongStore(t *testing.T) {
	newSong := songvote.Song{
		ID:     1,
		Name:   "Creep",
		Artist: "Radiohead",
	}
	storeUnderTest := songvote.NewInMemorySongStore

	t.Run("adding song returns valid ID", func(t *testing.T) {
		store := storeUnderTest()
		id, err := store.AddSong(newSong)
		assert.NoError(t, err)
		if id <= 0 {
			t.Errorf("got bad id %d", id)
		}
	})

	t.Run("can add and retreive a song", func(t *testing.T) {
		store := storeUnderTest()
		id, err := store.AddSong(newSong)
		assert.NoError(t, err)
		got, err := store.GetSong(id)
		assert.NoError(t, err)
		assert.Equal(t, got, newSong)
	})

	t.Run("can't add duplicate song", func(t *testing.T) {
		store := storeUnderTest()
		_, err := store.AddSong(newSong)
		assert.NoError(t, err)
		_, err = store.AddSong(newSong)
		assert.Error(t, err)
	})
}
