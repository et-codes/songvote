package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	s, _ := NewStore(":memory:")
	req := NewUserRequest{"John Doe", "password"}
	id, err := s.CreateUser(req)

	assert.NoError(t, err)
	assert.Equal(t, id, int64(1))
	t.Log(id)

	user, err := s.GetUser(id)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "password", user.Password)
	assert.False(t, user.Inactive)
	assert.Greater(t, user.Vetoes, 0)
}
