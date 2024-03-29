package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserStore(t *testing.T) {
	s, err := NewStore(":memory:")
	assert.NoError(t, err)

	t.Run("creates user and gets id", func(t *testing.T) {
		req := NewUserRequest{"John Doe", "password"}
		id, err := s.CreateUser(req)
		assert.NoError(t, err)
		assert.Equal(t, id, int64(1))
	})

	t.Run("created user contains correct info", func(t *testing.T) {
		user, err := s.GetUserByID(1)
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
		assert.NotEmpty(t, user.Password)
		assert.False(t, user.Inactive)
		assert.Greater(t, user.Vetoes, 0)
	})

	t.Run("get fails on non-existent user", func(t *testing.T) {
		_, err := s.GetUserByID(999)
		assert.Error(t, err)
	})

	t.Run("userExists works", func(t *testing.T) {
		exists := s.usernameExists("John Doe")
		assert.True(t, exists)

		exists = s.usernameExists("Aloysius Abercrombie")
		assert.False(t, exists)
	})

	t.Run("cannot create duplicate user", func(t *testing.T) {
		req := NewUserRequest{"John Doe", "password"}
		_, err := s.CreateUser(req)
		assert.Error(t, err)
	})

	t.Run("get user by name works", func(t *testing.T) {
		user, err := s.GetUserByName("John Doe")
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
	})

	t.Run("can get all users", func(t *testing.T) {
		users, err := s.GetUsers()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})

	t.Run("can delete a user", func(t *testing.T) {
		err := s.DeleteUser(1)
		assert.NoError(t, err)

		_, err = s.GetUserByID(1)
		assert.Error(t, err)
	})

	t.Run("can update a user", func(t *testing.T) {
		req := NewUserRequest{"John Blow", "password"}
		id, err := s.CreateUser(req)
		assert.NoError(t, err)

		updatedUser := User{
			ID:       id,
			Name:     "John Dough",
			Password: "new_password",
			Inactive: false,
			Vetoes:   2,
		}

		err = s.UpdateUser(&updatedUser)
		assert.NoError(t, err)

		user, err := s.GetUserByID(id)
		assert.NoError(t, err)
		assert.Equal(t, "John Dough", user.Name)
		assert.False(t, user.Inactive)
		assert.Equal(t, 2, user.Vetoes)
	})
}

func TestSongStore(t *testing.T) {

	var s *Store
	var user *User
	var err error

	t.Run("set up store and user", func(t *testing.T) {
		s, err = NewStore(":memory:")
		assert.NoError(t, err)

		req := NewUserRequest{"John Doe", "password"}
		_, err = s.CreateUser(req)
		assert.NoError(t, err)

		req = NewUserRequest{"Jane Doe", "password"}
		_, err = s.CreateUser(req)
		assert.NoError(t, err)

		user, err = s.GetUserByName("John Doe")
		assert.NoError(t, err)
	})

	t.Run("creates song and returns id", func(t *testing.T) {
		req := NewSongRequest{
			Title:   "Mirror In The Bathroom",
			Artist:  "Oingo Boingo",
			LinkURL: "https://youtu.be/SHWrmIzgB5A",
			AddedBy: user.ID,
		}
		id, err := s.CreateSong(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("can get song", func(t *testing.T) {
		song, err := s.GetSongByID(1)
		assert.NoError(t, err)
		assert.Equal(t, song.ID, int64(1))
		assert.Equal(t, song.Title, "Mirror In The Bathroom")
		assert.Equal(t, song.Artist, "Oingo Boingo")
		assert.Equal(t, song.LinkURL, "https://youtu.be/SHWrmIzgB5A")
		assert.Equal(t, song.AddedBy, int64(1))
		assert.Equal(t, song.Votes, 1)
	})

	t.Run("cannot create duplicate song/artist", func(t *testing.T) {
		req := NewSongRequest{
			Title:   "Mirror In The Bathroom",
			Artist:  "Oingo Boingo",
			LinkURL: "https://youtu.be/SHWrmIzgB5A",
			AddedBy: user.ID,
		}
		_, err := s.CreateSong(req)
		assert.Error(t, err)
	})

	t.Run("creating song records a vote by submitter", func(t *testing.T) {
		votes, err := s.GetVotesBySongID(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(votes))
		if len(votes) > 0 {
			assert.Equal(t, int64(1), votes[0].UserID)
		}
	})

	t.Run("can vote for a song", func(t *testing.T) {
		req := VoteRequest{1, 2}
		_, err := s.VoteForSong(req)
		assert.NoError(t, err)

		song, err := s.GetSongByID(1)
		assert.NoError(t, err)
		assert.Equal(t, 2, song.Votes)
	})

	t.Run("user cannot vote for same song twice", func(t *testing.T) {
		req := VoteRequest{1, 1}
		_, err := s.VoteForSong(req)
		assert.Error(t, err)
	})

	t.Run("can get all songs", func(t *testing.T) {
		songs, err := s.GetSongs()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(songs))

		song, _ := s.GetSongByID(1)
		assert.Equal(t, song, songs[0])
	})

	t.Run("voting for nonexistent song or user ID", func(t *testing.T) {
		req := VoteRequest{999, 1}
		_, err := s.VoteForSong(req)
		assert.Error(t, err)

		req = VoteRequest{1, 999}
		_, err = s.VoteForSong(req)
		assert.Error(t, err)
	})

	t.Run("can veto a song", func(t *testing.T) {
		req := VetoRequest{1, 1}
		id, err := s.VetoSong(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("cannot veto a song twice", func(t *testing.T) {
		req := VetoRequest{1, 2}
		_, err := s.VetoSong(req)
		assert.Error(t, err)
	})

	t.Run("user can't veto without vetoes", func(t *testing.T) {
		songReq := NewSongRequest{
			Title:   "Some Other Song",
			Artist:  "No Oingos Or Boingos",
			LinkURL: "https://youtu.be/SHWrmIzgB5A",
			AddedBy: user.ID,
		}
		id, err := s.CreateSong(songReq)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		vetoReq := VetoRequest{2, 1}
		_, err = s.VetoSong(vetoReq)
		assert.Error(t, err)
	})
}
