package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

// Unit tests for a Store.

var newSong = songvote.Song{
	ID:      1,
	Name:    "Fake Song #27",
	Artist:  "Fake Artist",
	LinkURL: "http://test.com",
	Votes:   12,
	Vetoed:  true,
}

func TestAddSongToStore(t *testing.T) {
	teardownSuite, store, _ := setupSuite(t)
	defer teardownSuite(t)

	t.Run("adding song returns valid ID", func(t *testing.T) {
		id, err := store.AddSong(newSong)
		assert.NoError(t, err)
		if id <= 0 {
			t.Errorf("got bad id %d", id)
		}
	})

	t.Run("can't add duplicate song", func(t *testing.T) {
		_, err := store.AddSong(newSong)
		assert.Error(t, err)
	})
}

func TestGetSongsFromStore(t *testing.T) {
	teardownSuite, store, _ := setupSuite(t)
	defer teardownSuite(t)

	t.Run("gets all songs in store", func(t *testing.T) {
		_, _ = store.AddSong(newSong)
		newSong.Name = "Fake Song #4"
		_, err := store.AddSong(newSong)
		assert.NoError(t, err)

		got := store.GetSongs()
		assert.Equal(t, len(got), 2)
	})
}

func TestDeleteSongFromStore(t *testing.T) {
	teardownSuite, store, _ := setupSuite(t)
	defer teardownSuite(t)

	t.Run("deletes a song", func(t *testing.T) {
		id, _ := store.AddSong(newSong)
		err := store.DeleteSong(id)
		assert.NoError(t, err)
	})
}

func TestUpdateSongInStore(t *testing.T) {
	teardownSuite, store, _ := setupSuite(t)
	defer teardownSuite(t)

	t.Run("updates song name", func(t *testing.T) {
		originalName := newSong.Name
		newName := "Fake Song #4"

		id, _ := store.AddSong(newSong)
		song, _ := store.GetSong(id)
		assert.Equal(t, song.Name, originalName)

		newSong.Name = newName
		err := store.UpdateSong(id, newSong)
		assert.NoError(t, err)

		song, _ = store.GetSong(id)
		assert.Equal(t, song.Name, newName)
	})
}

func TestAddVoteToSongInStore(t *testing.T) {
	teardownSuite, store, _ := setupSuite(t)
	defer teardownSuite(t)

	t.Run("updates vote count", func(t *testing.T) {
		id, _ := store.AddSong(newSong)
		err := store.AddVote(id)
		assert.NoError(t, err)

		song, _ := store.GetSong(id)
		assert.Equal(t, song.Votes, 13)
	})
}

func TestVetoSongInStore(t *testing.T) {
	teardownSuite, store, _ := setupSuite(t)
	defer teardownSuite(t)

	t.Run("sets veto value to true", func(t *testing.T) {
		id, _ := store.AddSong(newSong)
		err := store.Veto(id)
		assert.NoError(t, err)

		song, _ := store.GetSong(id)
		assert.True(t, song.Vetoed)
	})
}
