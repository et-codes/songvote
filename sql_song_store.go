package songvote

import (
	"context"
	"database/sql"
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

	return &SQLSongStore{
		db:  db,
		ctx: ctx,
	}
}

// GetSong returns a Song object with the given ID, or an error if it cannot
// be found.
func (s *SQLSongStore) GetSong(id int) (Song, error) {
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
func (s *SQLSongStore) AddSong(song Song) (int, error) {
	return 0, nil
}
