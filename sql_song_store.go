package songvote

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

const (
	dbDriver = "sqlite"
	dbPath   = "./db/songs.db"
)

// SQLSongStore is a song store backed by a SQL database.
type SQLSongStore struct {
	db  *sql.DB
	ctx context.Context
}

// NewSQLSongStore returns a pointer to a newly initialized store.
func NewSQLSongStore(dbPath string) *SQLSongStore {
	ctx := context.Background()
	db, err := sql.Open(dbDriver, dbPath)
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	log.Printf("Opened db %q.\n", dbPath)

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("error pinging db: %v", err)
	}
	log.Printf("Connected to db %q.\n", dbPath)

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
	row := s.db.QueryRowContext(s.ctx, "SELECT * FROM songs WHERE id = $1", id)
	song, err := rowToSong(row)

	switch {
	case err == sql.ErrNoRows:
		return song, fmt.Errorf("song ID %d not found", id)
	case err != nil:
		return song, fmt.Errorf("error getting song ID %d: %v", id, err)
	default:
		return song, nil
	}
}

// GetSongs returns a slice of Song objects representing all of the songs in
// the store.
func (s *SQLSongStore) GetSongs() Songs {
	rows, err := s.db.QueryContext(s.ctx, "SELECT * FROM songs")
	if err != nil {
		log.Fatalf("error querying songs from store: %v", err)
	}
	defer rows.Close()

	songs, err := rowsToSongs(rows)
	if err != nil {
		log.Fatal(err)
	}

	return songs
}

// AddSong adds the given Song object to the store, and returns the ID and
// an error if there was one.  Error will be returned when attempting to add
// a song with Name and Artist combination that is already in the store.
func (s *SQLSongStore) AddSong(song Song) (int64, error) {
	if s.songExists(song) {
		return 0, fmt.Errorf("%q by %q already exists", song.Name, song.Artist)
	}

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

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error retreiving new song ID: %v", err)
	}

	return id, nil
}

// DeleteSong will delete the given song ID from the table.
func (s *SQLSongStore) DeleteSong(id int64) error {
	_, err := s.db.ExecContext(s.ctx, "DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting song %d: %v", id, err)
	}
	return nil
}

func (s *SQLSongStore) UpdateSong(id int64, song Song) error {
	return fmt.Errorf("UpdateSong not implemented.")
}

func (s *SQLSongStore) AddVote(id int64) error {
	return fmt.Errorf("AddVote not implemented.")
}

func (s *SQLSongStore) AddVeto(id int64) error {
	return fmt.Errorf("AddVeto not implemented.")
}

// createSongsTable creates the database table for Songs if it does not
// already exist.
func (s *SQLSongStore) createSongsTable() error {
	_, err := s.db.ExecContext(s.ctx,
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

	return nil
}

// songExists queries the database for any song with the same name and artist
// as the given song and returns true if there is a match.
func (s *SQLSongStore) songExists(song Song) bool {
	var (
		name   string
		artist string
	)

	err := s.db.QueryRowContext(s.ctx,
		"SELECT name, artist FROM songs WHERE name = $1 AND artist = $2",
		song.Name,
		song.Artist,
	).Scan(&name, &artist)

	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Fatalf("error checking for song %q: %v\n", song.Name, err)
		return false
	default:
		return true
	}
}

// rowToSong marshals a *sql.Row result into a Song struct.
func rowToSong(row *sql.Row) (Song, error) {
	var song Song
	err := row.Scan(
		&song.ID,
		&song.Name,
		&song.Artist,
		&song.LinkURL,
		&song.Votes,
		&song.Vetoed,
	)
	return song, err
}

// rowsToSongs marshals a *sql.Rows result into a slice of Song structs.
func rowsToSongs(rows *sql.Rows) (Songs, error) {
	songs := Songs{}
	for rows.Next() {
		var song Song
		if err := rows.Scan(
			&song.ID,
			&song.Name,
			&song.Artist,
			&song.LinkURL,
			&song.Votes,
			&song.Vetoed,
		); err != nil {
			return songs, fmt.Errorf("problem scanning rows: %v", err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}
