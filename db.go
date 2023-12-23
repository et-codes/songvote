package main

import "fmt"

// CreateTables creates all tables required for the app.
func (s *Store) CreateTables() error {
	tableFuncs := map[string]func() error{
		"sessions": s.createSessionsTable,
		"users":    s.createUsersTable,
		"songs":    s.createSongsTable,
		"votes":    s.createVotesTable,
		"vetoes":   s.createVetoesTable,
	}

	for name, tf := range tableFuncs {
		if err := tf(); err != nil {
			return fmt.Errorf("error creating table %q: %v", name, err)
		}
	}

	return nil
}

// createSessionsTable creates the sessions table in the db if it doesn't exist.
func (s *Store) createSessionsTable() error {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			data BLOB NOT NULL,
			expiry REAL NOT NULL
		);
		CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry);`)
	return err
}

// createUsersTable creates the users table in the db if it doesn't exist.
func (s *Store) createUsersTable() error {
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

// createSongsTable creates the songs table in the db if it doesn't exist.
func (s *Store) createSongsTable() error {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS songs (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			artist TEXT NOT NULL,
			link_url TEXT,
			votes INTEGER,
			vetoed BOOLEAN,
			added_by INTEGER NOT NULL,
			FOREIGN KEY(added_by) REFERENCES users(id)
		);`)
	return err
}

// createVotesTable creates the votes table in the db if it doesn't exist.
func (s *Store) createVotesTable() error {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS votes (
			id INTEGER PRIMARY KEY,
			song_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY(song_id) REFERENCES songs(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`)
	return err
}

// createVetoesTable creates the vetoes table in the db if it doesn't exist.
func (s *Store) createVetoesTable() error {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS vetoes (
			id INTEGER PRIMARY KEY,
			song_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY(song_id) REFERENCES songs(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`)
	return err
}
