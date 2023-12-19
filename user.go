package main

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Inactive bool   `json:"inactive"`
	Vetoes   int    `json:"vetoes"`
}

type NewUserRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type NewUserResponse struct {
	ID int64 `json:"id"`
}

const (
	defaultVetoes = 1
)
