package songvote

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

const (
	driver = "sqlite"
	dbPath = "./db/songs.db"
)

// SQLSongStore is a song store backed by a SQL database.
type SQLSongStore struct {
	db  *sql.DB
	ctx context.Context
}

// NewSQLSongStore returns a pointer to a newly initialized store.
func NewSQLSongStore() *SQLSongStore {
	ctx := context.Background()
	db, err := sql.Open(driver, dbPath)
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("error pinging db: %v", err)
	}

	store := &SQLSongStore{
		db:  db,
		ctx: ctx,
	}

	if err := store.createSongsTable(); err != nil {
		log.Fatalf("error creating table: %v", err)
	}

	return store
}

// GetSong returns a Song object with the given ID, or an error if it cannot
// be found.
func (s *SQLSongStore) GetSong(id int64) (Song, error) {
	return Song{}, nil
}

// GetSongs returns a slice of Song objects representing all of the songs in
// the store.
func (s *SQLSongStore) GetSongs() []Song {
	return []Song{}
}

// AddSong adds the given Song object to the store, and returns the ID and
// an error if there was one.  Error will be returned when attempting to add
// a song with Name and Artist combination that is already in the store.
func (s *SQLSongStore) AddSong(song Song) (int64, error) {
	result, err := s.db.ExecContext(s.ctx,
		`INSERT INTO songs(name, artist, link_url, votes, vetoed) 
			VALUES ($1, $2, $3, $4, $5)`,
		song.Name,
		song.Artist,
		song.LinkURL,
		song.Votes,
		song.Vetoed,
	)
	if err != nil {
		return 0, fmt.Errorf("error inserting song into db: %v", err)
	}

	log.Printf("added song %s to the database", song.Name)

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error retreiving new song ID: %v", err)
	}

	return id, nil
}

// createSongsTable creates the database table for Songs if it does not
// already exist.
func (s *SQLSongStore) createSongsTable() error {
	result, err := s.db.ExecContext(s.ctx,
		`CREATE TABLE IF NOT EXISTS songs (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			artist TEXT NOT NULL,
			link_url TEXT,
			votes INTEGER,
			vetoed BOOLEAN
		)`,
	)

	if err != nil {
		return err
	}

	log.Printf("created user table in db: %v\n", result)

	return nil
}
