package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	s, _ := NewStore(":memory:")

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
		assert.Equal(t, "password", user.Password)
		assert.False(t, user.Inactive)
		assert.Greater(t, user.Vetoes, 0)
	})

	t.Run("get fails on non-existent user", func(t *testing.T) {
		_, err := s.GetUserByID(999)
		assert.Error(t, err)
	})

	t.Run("userExists works", func(t *testing.T) {
		exists := s.userExists("John Doe")
		assert.True(t, exists)

		exists = s.userExists("Aloysius Abercrombie")
		assert.False(t, exists)
	})

	t.Run("cannot create duplicate user", func(t *testing.T) {
		req := NewUserRequest{"John Doe", "password"}
		_, err := s.CreateUser(req)
		assert.Error(t, err)
	})
}
