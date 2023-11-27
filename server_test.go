// The purpose of these tests is to ensure that Server properly receives
// the HTTP requests, calls the appropriate methods, and returns the
// correct data types.

package songvote_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/et-codes/songvote"
	"github.com/et-codes/songvote/internal/assert"
)

func TestAddUserToServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("stores user when POST", func(t *testing.T) {
		request := newAddUserRequest(testUser)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusCreated)
		if len(store.AddUserCalls) == 0 {
			t.Fatal("call to AddUser not received")
		}

		var id int64
		err := songvote.UnmarshalJSON(response.Body, &id)
		assert.NoError(t, err)
		assert.Equal(t, id, int64(1))

		if !testUser.Equal(store.AddUserCalls[0]) {
			t.Errorf("did not store correct song, got %v, want %v",
				store.AddSongCalls[0], testUser)
		}
	})
}

func TestGetAllUsersFromServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("get all users", func(t *testing.T) {
		request := newGetUsersRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		users := songvote.Users{}
		_ = songvote.UnmarshalJSON(response.Body, &users)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, users, songvote.Users{})
		assert.Equal(t, store.GetUsersCallCount, 1)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/users", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestGetUserFromServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("get single user", func(t *testing.T) {
		request := newGetUserRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, len(store.GetUserCalls), 1)
	})

	t.Run("returns 404 if user ID not found", func(t *testing.T) {
		request := newGetUserRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestUpdateUserOnServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	populateWithUser(server, testUser)

	t.Run("update user", func(t *testing.T) {
		newUser := testUser
		newUser.Name = "Changed Name"
		request := newUpdateUserRequest(1, newUser)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, len(store.UpdateUserCalls), 1)
		if !newUser.Equal(store.UpdateUserCalls[1]) {
			t.Errorf("got %v, want %v", store.UpdateUserCalls[1], newUser)
		}
	})

	t.Run("returns error if user not found", func(t *testing.T) {
		request := newUpdateUserRequest(10, testUser)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}

func TestGetAllSongsFromServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("get all songs", func(t *testing.T) {
		request := newGetSongsRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		songs := songvote.Songs{}
		_ = songvote.UnmarshalJSON(response.Body, &songs)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, songs, songvote.Songs{})
		assert.Equal(t, store.GetSongsCallCount, 1)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/songs", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestGetSongFromServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("get single song", func(t *testing.T) {
		request := newGetSongRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, len(store.GetSongCalls), 1)
	})

	t.Run("returns 404 if song ID not found", func(t *testing.T) {
		request := newGetSongRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/songs/0", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestAddSongsToServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("stores song when POST", func(t *testing.T) {
		newSong := songvote.Song{
			Name:   "Creep",
			Artist: "Radiohead",
		}

		request := newAddSongRequest(newSong)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusCreated)
		if len(store.AddSongCalls) == 0 {
			t.Fatal("call to AddSong not received")
		}

		var id int64
		_ = songvote.UnmarshalJSON(response.Body, &id)
		assert.Equal(t, id, int64(1))

		if !newSong.Equal(store.AddSongCalls[0]) {
			t.Errorf("did not store correct song, got %v, want %v",
				store.AddSongCalls[0], newSong)
		}
	})
}

func TestDeleteUserFromServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("delete user response", func(t *testing.T) {
		request := newDeleteUserRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, store.DeleteUserCalls[0], int64(1))
	})

	t.Run("returns error if user not found", func(t *testing.T) {
		request := newDeleteUserRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
		assert.Equal(t, store.DeleteUserCalls[1], int64(10))
	})
}

func TestDeleteSongFromServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("delete song response", func(t *testing.T) {
		request := newDeleteSongRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, store.DeleteSongCalls[0], int64(1))
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		request := newDeleteSongRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
		assert.Equal(t, store.DeleteSongCalls[1], int64(10))
	})
}

func TestUpdateSongOnServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	song := songvote.Song{
		Name:   "Creep",
		Artist: "Radiohead",
	}

	t.Run("update song", func(t *testing.T) {
		request := newUpdateSongRequest(1, song)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, len(store.UpdateSongCalls), 1)
		if !song.Equal(store.UpdateSongCalls[1]) {
			t.Errorf("got %v, want %v", store.UpdateSongCalls[1], song)
		}
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		request := newUpdateSongRequest(10, song)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})
}

func TestAddVoteOnServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("update vote count", func(t *testing.T) {
		request := newVoteRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, store.AddVoteCalls[0], int64(1))
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/songs/vote/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		request := newVoteRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusInternalServerError)
	})
}

func TestVetoSongOnServer(t *testing.T) {
	teardownSuite, store, server := setupStubSuite(t)
	defer teardownSuite(t)

	t.Run("veto a song", func(t *testing.T) {
		request := newVetoRequest(1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNoContent)
		assert.Equal(t, len(store.VetoCalls), 1)
	})

	t.Run("returns error if song not found", func(t *testing.T) {
		request := newVetoRequest(10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusInternalServerError)
	})

	t.Run("returns 405 when wrong method used", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/songs/veto/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}
