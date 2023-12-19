package main

import "math/rand"

type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"username"`
	Password   string `json:"password"`
	IsInactive bool   `json:"isInactive"`
	Vetoes     int    `json:"vetoes"`
}

type NewUserRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

const (
	defaultVetoes = 1
)

func AddUser(req NewUserRequest) *User {
	id := rand.Uint64()
	return &User{
		ID:       id,
		Name:     req.Name,
		Password: req.Password,
		Vetoes:   defaultVetoes,
	}
}
