package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	req := NewUserRequest{"John Doe", "password"}
	got := NewUser(req)

	assert.Greater(t, got.ID, uint64(0))
	assert.Equal(t, "John Doe", got.Name)
	assert.Equal(t, "password", got.Password)
	assert.False(t, got.IsInactive)
	assert.Greater(t, got.Vetoes, 0)
}
