package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

const (
	memoryDB   = ":memory:"              // in-memory database
	testFileDB = "./db/songvote_test.db" // persistent test database
)

// Configure the store to use either the in-memory db or disk-based db. Use
// the disk-based db if you want to populate it with data for more extensive
// testing.
var store = songvote.NewSQLStore(memoryDB)

var newSong = songvote.Song{
	ID:      1,
	Name:    "Fake Song #27",
	Artist:  "Fake Artist",
	LinkURL: "http://test.com",
	Votes:   12,
	Vetoed:  true,
}

func TestAddSong(t *testing.T) {
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
}

func TestGetSongs(t *testing.T) {
	t.Run("gets all songs in store", func(t *testing.T) {
		_, _ = addSong(t, newSong)
		newSong.Name = "Fake Song #4"
		_, err := addSong(t, newSong)
		assert.NoError(t, err)

		got := store.GetSongs()
		assert.Equal(t, len(got), 2)
	})
}

func TestDeleteSong(t *testing.T) {
	t.Run("deletes a song", func(t *testing.T) {
		id, _ := addSong(t, newSong)
		err := store.DeleteSong(id)
		assert.NoError(t, err)
	})
}

func TestUpdateSong(t *testing.T) {
	t.Run("updates song name", func(t *testing.T) {
		originalName := newSong.Name
		newName := "Fake Song #4"

		id, _ := addSong(t, newSong)
		song, _ := store.GetSong(id)
		assert.Equal(t, song.Name, originalName)

		newSong.Name = newName
		err := store.UpdateSong(id, newSong)
		assert.NoError(t, err)

		song, _ = store.GetSong(id)
		assert.Equal(t, song.Name, newName)
	})
}

func TestAddVote(t *testing.T) {
	t.Run("updates vote count", func(t *testing.T) {
		id, _ := addSong(t, newSong)
		err := store.AddVote(id)
		assert.NoError(t, err)

		song, _ := store.GetSong(id)
		assert.Equal(t, song.Votes, 13)
	})
}

func TestVeto(t *testing.T) {
	t.Run("sets veto value to true", func(t *testing.T) {
		id, _ := addSong(t, newSong)
		err := store.Veto(id)
		assert.NoError(t, err)

		song, _ := store.GetSong(id)
		assert.True(t, song.Vetoed)
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
