package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// Store contains data related to storage.
type Store struct {
	db *sql.DB
}

// NewStore creates a new SQLite3 database store.
func NewStore(dbPath string) (*Store, error) {
	slog.Info("Opening db", "path", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error opening db: %v", err)
	}
	slog.Info("Connected to db.")

	store := &Store{db}

	if err := store.CreateTables(); err != nil {
		return nil, fmt.Errorf("error creating tables: %v", err)
	}

	return store, nil
}

// CreateUser creates a new user with the given request data.
func (s *Store) CreateUser(req NewUserRequest) (int64, error) {
	if s.usernameExists(req.Name) {
		return 0, ErrConflict
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, NewServerError(http.StatusInternalServerError, err.Error())
	}

	result, err := s.db.Exec(
		`INSERT INTO users(name, password, inactive, vetoes) VALUES($1, $2, $3, $4)`,
		req.Name, pwd, false, initialVetoes,
	)
	if err != nil {
		return 0, NewServerError(http.StatusInternalServerError, err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return id, NewServerError(http.StatusInternalServerError, err.Error())
	}

	slog.Info("New user created", "id", id, "name", req.Name)
	return id, nil
}

// GetUserByID returns user data that matches the given ID.
func (s *Store) GetUserByID(id int64) (*User, error) {
	user := User{}

	row := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Inactive, &user.Vetoes)
	if err != nil {
		return nil, ErrNotFound
	}

	return &user, nil
}

// GetUserByName returns user data that matches the given username.
func (s *Store) GetUserByName(username string) (*User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE name = $1", username)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Inactive, &user.Vetoes)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// usernameExists returns true if a user with the given name is in the database.
func (s *Store) usernameExists(username string) bool {
	row := s.db.QueryRow("SELECT id FROM users WHERE name = $1", username)
	var id int64
	if err := row.Scan(&id); err != nil {
		return false
	}
	return true
}

// userIDExists returns true if a user with the given ID is in the database.
func (s *Store) userIDExists(id int64) bool {
	row := s.db.QueryRow("SELECT id FROM users WHERE id = $1", id)
	var userID int64
	err := row.Scan(&userID)
	return err == nil
}

// CreateSong creates a new user with the given request data.
func (s *Store) CreateSong(req NewSongRequest) (int64, error) {
	if s.songTitleArtistExists(req.Title, req.Artist) {
		return 0, ErrConflict
	}

	result, err := s.db.Exec(
		`INSERT INTO songs(title, artist, link_url, votes, vetoed, added_by) 
		VALUES($1, $2, $3, $4, $5, $6)`,
		req.Title, req.Artist, req.LinkURL, 0, false, req.AddedBy,
	)
	if err != nil {
		return 0, NewServerError(http.StatusInternalServerError, err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return id, NewServerError(http.StatusInternalServerError, err.Error())
	}

	slog.Info("New song created", "id", id, "title", req.Title, "artist", req.Artist)

	voteReq := NewVoteRequest{SongID: id, UserID: req.AddedBy}
	_, err = s.VoteForSong(voteReq)
	if err != nil {
		return id, err
	}

	return id, nil
}

// GetSongByID returns song data that matches the given ID.
func (s *Store) GetSongByID(id int64) (*Song, error) {
	song := Song{}

	row := s.db.QueryRow("SELECT * FROM songs WHERE id = $1", id)
	err := row.Scan(&song.ID, &song.Title, &song.Artist, &song.LinkURL,
		&song.Votes, &song.Vetoed, &song.AddedBy)
	if err != nil {
		slog.Error("error retreiving song", "error", err)
		return nil, ErrNotFound
	}

	return &song, nil
}

// GetSongs returns all songs in the database.
func (s *Store) GetSongs() ([]*Song, error) {
	songs := []*Song{}

	rows, err := s.db.Query("SELECT * FROM songs")
	if err != nil {
		slog.Error("error getting songs from db", "error", err)
		return nil, err
	}

	for rows.Next() {
		song := Song{}
		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.LinkURL,
			&song.Votes, &song.Vetoed, &song.AddedBy)
		if err != nil {
			slog.Error("error scanning rows", "error", err)
		}
		songs = append(songs, &song)
	}

	return songs, nil
}

// songTitleArtistExists checks whether a title/artist combination already exists.
func (s *Store) songTitleArtistExists(title, artist string) bool {
	var id int64
	row := s.db.QueryRow("SELECT id FROM songs WHERE title = $1 AND artist = $2",
		title, artist)
	err := row.Scan(&id)
	return err == nil
}

// songIDExists returns true if a song with the given ID is in the database.
func (s *Store) songIDExists(id int64) bool {
	var songID int64
	row := s.db.QueryRow("SELECT id FROM songs WHERE id = $1", id)
	err := row.Scan(&songID)
	return err == nil
}

// GetVotesBySongID returns a slice of votes for the given song ID.
func (s *Store) GetVotesBySongID(songID int64) ([]Vote, error) {
	votes := []Vote{}
	rows, err := s.db.Query("SELECT * FROM votes WHERE song_id = $1", songID)
	if err != nil {
		slog.Error("error querying votes", "error", err)
		return nil, fmt.Errorf("error querying votes: %v", err)
	}

	for rows.Next() {
		vote := Vote{}
		err := rows.Scan(&vote.ID, &vote.SongID, &vote.UserID)
		if err != nil {
			slog.Error("Error scanning rows", "error", err)
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}
		votes = append(votes, vote)
	}

	return votes, nil
}

// VoteForSong adds a vote to a song.
func (s *Store) VoteForSong(req NewVoteRequest) (int64, error) {
	// Validate input.
	if req.SongID < 1 || req.UserID < 1 {
		return 0, fmt.Errorf("invalid song/user ID")
	}

	if !s.userIDExists(req.UserID) {
		return 0, fmt.Errorf("user %d not found", req.UserID)
	}

	if !s.songIDExists(req.SongID) {
		return 0, fmt.Errorf("song %d not found", req.SongID)
	}

	// Get existing votes for the song.
	votes, err := s.GetVotesBySongID(req.SongID)
	if err != nil {
		slog.Error("error getting votes", "error", err)
		return 0, fmt.Errorf("error getting votes: %v", err)
	}

	// Check if user has already voted for the song.
	for _, vote := range votes {
		if vote.SongID == req.SongID && vote.UserID == req.UserID {
			return 0, fmt.Errorf("user has already voted for this song")
		}
	}

	// Create the new vote record.
	id, err := s.createVote(req)
	if err != nil {
		return 0, err
	}
	slog.Info("New vote created", "id", id, "song_id", req.SongID, "user_id", req.UserID)

	// Update vote count on the song.
	_, err = s.db.Exec("UPDATE songs SET votes = $1 WHERE id = $2",
		len(votes)+1, req.SongID)
	if err != nil {
		slog.Error("error updating vote count", "error", err)
		return id, fmt.Errorf("error updating vote count: %v", err)
	}

	return id, nil
}

// createVote adds a vote record to the database.
func (s *Store) createVote(req NewVoteRequest) (int64, error) {
	result, err := s.db.Exec("INSERT INTO votes(song_id, user_id) VALUES($1, $2)",
		req.SongID, req.UserID)
	if err != nil {
		slog.Error("Error recording vote", "error", err)
		return 0, fmt.Errorf("error recording vote: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("error retreiving vote id", "error", err)
		return id, fmt.Errorf("error retreiving vote id: %v", err)
	}

	return id, nil
}
