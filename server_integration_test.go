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
