package songvote_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

// These integration tests verify function of a SongStore through the API.

func TestAddUsers(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	t.Run("add user and get ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newAddUserRequest(testUser)

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "1"

		assert.Equal(t, response.Code, http.StatusCreated)
		assert.Equal(t, got, want)
	})

	t.Run("returns 409 with duplicate user", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newAddUserRequest(testUser)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusConflict)
	})
}

func TestGetUser(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithUsers(server, userTestDataFile)

	t.Run("get all users", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetUsersRequest()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)

		users := songvote.Users{}
		err := songvote.UnmarshalJSON(response.Body, &users)
		assert.NoError(t, err)

		assert.Equal(t, len(users), 4)
		if !users[0].Equal(testUser) {
			t.Errorf("want %v, got %v", testUser, users[0])
		}
	})
}

func TestDeleteUser(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithUser(server, testUser)

	t.Run("delete a user and cannot retreive it", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newDeleteUserRequest(1)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)

		response = httptest.NewRecorder()
		request = newGetUserRequest(1)
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns error with unknown ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newDeleteUserRequest(999)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}

func TestUpdateUser(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithUser(server, testUser)

	t.Run("user is updated", func(t *testing.T) {
		newName := "Test Name"
		renamedUser := testUser
		renamedUser.Name = newName

		response := httptest.NewRecorder()
		request := newUpdateUserRequest(1, renamedUser)

		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusNoContent)

		response = httptest.NewRecorder()
		request = newGetUserRequest(1)

		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		updatedUser := songvote.User{}
		err := songvote.UnmarshalJSON(response.Body, &updatedUser)
		assert.NoError(t, err)
		assert.Equal(t, updatedUser.Name, newName)
	})
}

func TestAddSongs(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	t.Run("add song and get ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newAddSongRequest(testSong)

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "1"

		assert.Equal(t, response.Code, http.StatusCreated)
		assert.Equal(t, got, want)
	})

	t.Run("returns 409 with duplicate song", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newAddSongRequest(testSong)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusConflict)
	})
}

func TestGetSong(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSongs(server, songTestDataFile)

	t.Run("get song with ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetSongRequest(1)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)

		got := songvote.Song{}
		err := songvote.UnmarshalJSON(response.Body, &got)
		assert.NoError(t, err)

		want := testSong
		want.ID = 1

		assert.Equal(t, got, want)
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetSongRequest(999)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns error if invalid ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/songs/thatonesong", nil)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusBadRequest)
	})

	t.Run("get all songs", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetSongsRequest()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)

		songs := songvote.Songs{}
		err := songvote.UnmarshalJSON(response.Body, &songs)
		assert.NoError(t, err)

		assert.Equal(t, len(songs), 5)
		if !songs[0].Equal(testSong) {
			t.Errorf("want %v, got %v", testSong, songs[0])
		}
	})
}

func TestDeleteSong(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("delete a song and cannot retreive it", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newDeleteSongRequest(1)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)

		response = httptest.NewRecorder()
		request = newGetSongRequest(1)
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns error with unknown ID", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newDeleteSongRequest(999)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}

func TestUpdateSong(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithSong(server, testSong)

	t.Run("song is updated", func(t *testing.T) {
		newName := "Test Name"
		renamedSong := testSong
		renamedSong.Name = newName

		response := httptest.NewRecorder()
		request := newUpdateSongRequest(1, renamedSong)

		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusNoContent)

		response = httptest.NewRecorder()
		request = newGetSongRequest(1)

		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		updatedSong := songvote.Song{}
		err := songvote.UnmarshalJSON(response.Body, &updatedSong)
		assert.NoError(t, err)
		assert.Equal(t, updatedSong.Name, newName)
	})
}

func TestVoteForSong(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithUsers(server, userTestDataFile)
	populateWithSongs(server, songTestDataFile)

	vote1 := songvote.Vote{SongID: 1, UserID: 1}
	vote2 := songvote.Vote{SongID: 1, UserID: 2}

	t.Run("vote count increments", func(t *testing.T) {
		// Make vote.
		response := httptest.NewRecorder()
		request := newVoteRequest(vote1)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)

		// Check song's vote count.
		response = httptest.NewRecorder()
		request = newGetSongRequest(vote1.UserID)

		server.ServeHTTP(response, request)

		updatedSong := songvote.Song{}
		err := songvote.UnmarshalJSON(response.Body, &updatedSong)
		assert.NoError(t, err)
		assert.Equal(t, updatedSong.Votes, 11)
	})

	t.Run("only one vote per user per song", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newVoteRequest(vote1)
		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusConflict)
	})

	t.Run("can retreive list of votes", func(t *testing.T) {
		// Add another vote to make it interesting.
		response := httptest.NewRecorder()
		request := newVoteRequest(vote2)
		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)

		// Get votes for song.
		response = httptest.NewRecorder()
		request = newGetVotesRequest(1)
		server.ServeHTTP(response, request)

		votes := songvote.Votes{}
		err := songvote.UnmarshalJSON(response.Body, &votes)
		assert.NoError(t, err)
		assert.Equal(t, len(votes), 2)
		for _, vote := range votes {
			if vote.UserID != int64(1) && vote.UserID != int64(2) {
				t.Errorf("unexpected voter %d", vote.UserID)
			}
		}
	})
}

func TestVetoSong(t *testing.T) {
	teardownSuite, _, server := setupSuite(t)
	defer teardownSuite(t)

	populateWithUsers(server, userTestDataFile)
	populateWithSongs(server, songTestDataFile)

	veto := songvote.Veto{1, 1, 1}
	veto2 := songvote.Veto{1, 1, 2}

	t.Run("sets veto to true", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newVetoRequest(veto)
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusNoContent)

		response = httptest.NewRecorder()
		request = newGetSongRequest(veto.SongID)
		server.ServeHTTP(response, request)
		song := songvote.Song{}
		err := songvote.UnmarshalJSON(response.Body, &song)
		assert.NoError(t, err)
		assert.True(t, song.Vetoed)
	})

	t.Run("no error when vetoing twice", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newVetoRequest(veto2)
		server.ServeHTTP(response, request)
		if response.Code != http.StatusNoContent {
			serverError := songvote.ServerError{}
			_ = songvote.UnmarshalJSON(response.Body, &serverError)
			t.Errorf("didn't want error, got %d: %+v", response.Code, serverError)
		}
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		veto.SongID = 99
		request := newVetoRequest(veto)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/songs/veto", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}
