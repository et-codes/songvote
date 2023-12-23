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
	if s.userExists(req.Name) {
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

// userExists returns true if a user with the given name is in the database.
func (s *Store) userExists(username string) bool {
	row := s.db.QueryRow("SELECT id FROM users WHERE name = $1", username)
	var id int64
	if err := row.Scan(&id); err != nil {
		return false
	}
	return true
}

// CreateSong creates a new user with the given request data.
func (s *Store) CreateSong(req NewSongRequest) (int64, error) {
	// if s.songExists(req.Title, req.Artist) {
	// 	return 0, ErrConflict
	// }

	result, err := s.db.Exec(
		`INSERT INTO songs(title, artist, link_url, votes, vetoed, added_by) 
		VALUES($1, $2, $3, $4, $5, $6)`,
		req.Title, req.Artist, req.LinkURL, 1, false, req.AddedBy,
	)
	if err != nil {
		return 0, NewServerError(http.StatusInternalServerError, err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return id, NewServerError(http.StatusInternalServerError, err.Error())
	}

	slog.Info("New song created", "id", id, "title", req.Title, "artist", req.Artist)
	return id, nil
}
