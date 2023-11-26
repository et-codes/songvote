package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote/internal/assert"
)

func TestAddSongToStore(t *testing.T) {
	teardownSuite, store, _ := setupSuite(t)
	defer teardownSuite(t)

	t.Run("adding song returns valid ID", func(t *testing.T) {
		id, err := store.AddSong(testSong)
		assert.NoError(t, err)
		if id <= 0 {
			t.Errorf("got bad id %d", id)
		}
	})

	t.Run("can't add duplicate song", func(t *testing.T) {
		_, err := store.AddSong(testSong)
		assert.Error(t, err)
	})
}

func TestGetSongsFromStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("gets all songs in store", func(t *testing.T) {
		newSong := testSong
		newSong.Name = "Fake Song #4"
		_, err := store.AddSong(newSong)
		assert.NoError(t, err)

		got := store.GetSongs()
		assert.Equal(t, len(got), 2)
	})
}

func TestDeleteSongFromStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("deletes a song", func(t *testing.T) {
		err := store.DeleteSong(1)
		assert.NoError(t, err)
	})
}

func TestUpdateSongInStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("updates song name", func(t *testing.T) {
		newSong := testSong
		newSong.Name = "Fake Song #4"

		err := store.UpdateSong(1, newSong)
		assert.NoError(t, err)

		song, _ := store.GetSong(1)
		assert.Equal(t, song.Name, newSong.Name)
	})
}

func TestAddVoteToSongInStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("updates vote count", func(t *testing.T) {
		err := store.AddVote(1)
		assert.NoError(t, err)

		song, _ := store.GetSong(1)
		assert.Equal(t, song.Votes, 11)
	})
}

func TestVetoSongInStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("sets veto value to true", func(t *testing.T) {
		err := store.Veto(1)
		assert.NoError(t, err)

		song, _ := store.GetSong(1)
		assert.True(t, song.Vetoed)
	})
}
