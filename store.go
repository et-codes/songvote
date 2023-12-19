package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

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

	if err := store.createUserTable(); err != nil {
		return nil, fmt.Errorf("error creating user table: %v", err)
	}

	return store, nil
}

func (s *Store) CreateUser(req NewUserRequest) (int64, error) {
	if s.userExists(req.Name) {
		return 0, ErrConflict
	}

	result, err := s.db.Exec(
		`INSERT INTO users(name, password, inactive, vetoes) VALUES($1, $2, $3, $4)`,
		req.Name, req.Password, false, defaultVetoes,
	)
	if err != nil {
		return 0, NewServerError(http.StatusInternalServerError, err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return id, NewServerError(http.StatusInternalServerError, err.Error())
	}

	slog.Info("New user created", "id", id)
	return id, nil
}

func (s *Store) GetUserByID(id int64) (*User, error) {
	user := User{}

	row := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Inactive, &user.Vetoes)
	if err != nil {
		return nil, ErrNotFound
	}

	return &user, nil
}

func (s *Store) userExists(username string) bool {
	row := s.db.QueryRow("SELECT id FROM users WHERE name = $1", username)
	var id int64
	if err := row.Scan(&id); err != nil {
		return false
	}
	return true
}

func (s *Store) createUserTable() error {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			inactive BOOLEAN,
			vetoes INTEGER
		);`)
	return err
}
