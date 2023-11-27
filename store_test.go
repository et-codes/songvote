package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote/internal/assert"
)

func TestAddUserToStore(t *testing.T) {
	teardownSuite, store, _ := setupSuite(t)
	defer teardownSuite(t)

	t.Run("adding user returns valid ID", func(t *testing.T) {
		id, err := store.AddUser(testUser)
		assert.NoError(t, err)
		if id < 1 {
			t.Errorf("got bad id %d", id)
		}
	})
}

func TestGetAllUsersFromStore(t *testing.T) {
	t.Run("gets a list of all users", func(t *testing.T) {
		teardownSuite, store, server := setupSuite(t)
		defer teardownSuite(t)

		populateWithUsers(server, userTestDataFile)

		users := store.GetUsers()
		assert.Equal(t, len(users), 4)
	})

	t.Run("returns empty if no users in store", func(t *testing.T) {
		teardownSuite, store, _ := setupSuite(t)
		defer teardownSuite(t)

		users := store.GetUsers()
		assert.Equal(t, len(users), 0)
	})
}

func TestGetUserFromStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithUser(server, testUser)

	t.Run("gets user from store", func(t *testing.T) {
		got, err := store.GetUser(1)
		assert.NoError(t, err)
		if !got.Equal(testUser) {
			t.Errorf("got %v, want %v", got, testUser)
		}
	})

	t.Run("returns error if user not found", func(t *testing.T) {
		_, err := store.GetUser(999)
		assert.Error(t, err)
	})
}

func TestAddSongToStore(t *testing.T) {
	teardownSuite, store, _ := setupSuite(t)
	defer teardownSuite(t)

	t.Run("adding song returns valid ID", func(t *testing.T) {
		id, err := store.AddSong(testSong)
		assert.NoError(t, err)
		if id < 1 {
			t.Errorf("got bad id %d", id)
		}
	})

	t.Run("can't add duplicate song", func(t *testing.T) {
		_, err := store.AddSong(testSong)
		assert.Error(t, err)
	})
}

func TestGetSongsFromStore(t *testing.T) {
	t.Run("gets all songs in store", func(t *testing.T) {
		teardownSuite, store, server := setupSuite(t)
		defer teardownSuite(t)

		populateWithSongs(server, songTestDataFile)

		got := store.GetSongs()
		assert.Equal(t, len(got), 5)
	})

	t.Run("returns empty if no songs in store", func(t *testing.T) {
		teardownSuite, store, _ := setupSuite(t)
		defer teardownSuite(t)

		songs := store.GetSongs()
		assert.Equal(t, len(songs), 0)
	})
}

func TestGetSongFromStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("gets song from store", func(t *testing.T) {
		got, err := store.GetSong(1)
		assert.NoError(t, err)
		if !got.Equal(testSong) {
			t.Errorf("got %v, want %v", got, testSong)
		}
	})
}

func TestDeleteSongFromStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("deletes a song", func(t *testing.T) {
		err := store.DeleteSong(1)
		assert.NoError(t, err)
		_, err = store.GetSong(1)
		assert.Error(t, err)
	})
}

func TestUpdateSongInStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("updates song name", func(t *testing.T) {
		newSong := testSong
		newSong.Name = "Fake Song"

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
