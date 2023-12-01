package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote"
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

func TestDeleteUserFromStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithUser(server, testUser)

	t.Run("deletes a user", func(t *testing.T) {
		err := store.DeleteUser(1)
		assert.NoError(t, err)

		_, err = store.GetUser(1)
		assert.Error(t, err)

		err = store.DeleteUser(1)
		assert.Error(t, err)
	})
}

func TestUpdateUserInStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithUser(server, testUser)

	t.Run("updates user name", func(t *testing.T) {
		newUserData := songvote.User{
			Inactive: false,
			Name:     "Fake User",
			Password: "p3nc1l",
		}

		err := store.UpdateUser(1, newUserData)
		assert.NoError(t, err)

		user, _ := store.GetUser(1)
		assert.False(t, user.Inactive)
		assert.Equal(t, user.Name, newUserData.Name)
		assert.Equal(t, user.Password, newUserData.Password)
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

		err = store.DeleteSong(1)
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
	populateWithUser(server, testUser)

	t.Run("updates vote count", func(t *testing.T) {
		err := store.AddVote(songvote.Vote{
			SongID: 1,
			UserID: 1,
		})
		assert.NoError(t, err)

		song, _ := store.GetSong(1)
		assert.Equal(t, song.Votes, 11)
	})

	t.Run("inactive user cannot vote", func(t *testing.T) {
		inactiveUser := songvote.User{
			ID:       2,
			Inactive: true,
			Name:     "Jane Doe",
			Password: "p@ssword",
			Vetoes:   1,
		}

		populateWithUser(server, inactiveUser)

		err := store.AddVote(songvote.Vote{
			SongID: 1,
			UserID: 2,
		})
		assert.Error(t, err)
	})

	t.Run("tracks who voted for the song", func(t *testing.T) {
		want := songvote.Votes{
			{ID: 1, SongID: 1, UserID: 1},
		}
		got, err := store.GetVotesForSong(1)
		assert.NoError(t, err)
		assert.Equal(t, got, want)
	})

	t.Run("only one vote per user per song", func(t *testing.T) {
		err := store.AddVote(songvote.Vote{
			SongID: 1,
			UserID: 1,
		})
		assert.Error(t, err)

		got, err := store.GetVotesForSong(1)
		assert.NoError(t, err)
		assert.Equal(t, len(got), 1)
	})
}

func TestVetoSongInStore(t *testing.T) {
	teardownSuite, store, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSongs(server, songTestDataFile)
	populateWithUsers(server, userTestDataFile)

	veto := songvote.Veto{1, 1, 1}

	t.Run("sets veto value to true", func(t *testing.T) {
		err := store.Veto(veto)
		assert.NoError(t, err)

		song, _ := store.GetSong(1)
		assert.True(t, song.Vetoed)
	})

	t.Run("records who vetoed what", func(t *testing.T) {
		user, err := store.GetVetoedBy(1)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, testUser.Name)
	})

	t.Run("can't veto if user has no vetoes left", func(t *testing.T) {
		err := store.Veto(veto)
		assert.Error(t, err)
	})

	t.Run("returns error if song or user not found", func(t *testing.T) {
		err := store.Veto(songvote.Veto{1, 99, 2})
		assert.Error(t, err)

		err = store.Veto(songvote.Veto{1, 2, 99})
		assert.Error(t, err)
	})
}
