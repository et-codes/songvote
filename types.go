package main

const initialVetoes = 1

// User types

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Inactive bool   `json:"inactive"`
	Vetoes   int    `json:"vetoes"`
}

type NewUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type NewUserResponse struct {
	ID int64 `json:"id"`
}

// Song types

type Song struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	LinkURL string `json:"link_url"`
	Votes   int    `json:"votes"`
	Vetoed  bool   `json:"vetoed"`
	AddedBy int64  `json:"added_by"`
}

type NewSongRequest struct {
	AddedBy int64  `json:"added_by"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	LinkURL string `json:"link_url"`
}

// Vote types

type Vote struct {
	ID     int64 `json:"id"`
	SongID int64 `json:"song_id"`
	UserID int64 `json:"user_id"`
}

type NewVoteRequest struct {
	SongID int64 `json:"song_id"`
	UserID int64 `json:"user_id"`
}

// Veto types

type Veto struct {
	ID     int64 `json:"id"`
	SongID int64 `json:"song_id"`
	UserID int64 `json:"user_id"`
}

type NewVetoRequest struct {
	SongID int64 `json:"song_id"`
	UserID int64 `json:"user_id"`
}
