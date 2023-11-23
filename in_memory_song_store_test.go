package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote"
)

func TestInMemorySongStore(t *testing.T) {
	store := songvote.InMemorySongStore{}

	t.Run("can add and retreive a song", func(t *testing.T) {
		newSong := songvote.Song{
			Name:   "Creep",
			Artist: "Radiohead",
		}
		store.AddSong(newSong)
		got := store.GetSong(0)
		assertEqual(t, got, newSong)
	})

	t.Run("can't add duplicate song", func(t *testing.T) {
		newSong := songvote.Song{
			Name:   "Creep",
			Artist: "Radiohead",
		}
		store.AddSong(newSong)
		store.AddSong(newSong)

		got := store.GetSong(1)
		assertEqual(t, got, songvote.Song{})
	})
}
