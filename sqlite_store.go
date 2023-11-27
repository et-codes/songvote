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
)

// SQLiteStore is a data store backed by a SQL database.
type SQLiteStore struct {
	db  *sql.DB
	ctx context.Context
}

// NewSQLiteStore returns a pointer to a newly initialized store.
func NewSQLiteStore(dbPath string) *SQLiteStore {
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

	store := &SQLiteStore{
		db:  db,
		ctx: ctx,
	}

	if err := store.createTables(); err != nil {
		log.Fatalf("error creating table: %v", err)
	}

	return store
}

// AddUser adds the given User to the store.
func (s *SQLiteStore) AddUser(user User) (int64, error) {
	if s.userExists(user) {
		return 0, fmt.Errorf("user %q already exists", user.Name)
	}

	result, err := s.db.ExecContext(s.ctx,
		`INSERT INTO users(active, name, password, vetoes) 
			VALUES ($1, $2, $3, $4)`,
		true,
		user.Name,
		user.Password,
		user.Vetoes,
	)
	if err != nil {
		return 0, fmt.Errorf("error adding user to db: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error retreiving new user ID: %v", err)
	}

	return id, nil
}

// GetSong returns a Song object with the given ID, or an error if it cannot
// be found.
func (s *SQLiteStore) GetSong(id int64) (Song, error) {
	row := s.db.QueryRowContext(s.ctx,
		"SELECT * FROM songs WHERE id = $1",
		id,
	)
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
func (s *SQLiteStore) GetSongs() Songs {
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
func (s *SQLiteStore) AddSong(song Song) (int64, error) {
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
func (s *SQLiteStore) DeleteSong(id int64) error {
	result, err := s.db.ExecContext(s.ctx,
		"DELETE FROM songs WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting song %d: %v", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("song %d not found: %v", id, err)
	}

	return nil
}

// UpdateSong allows changing the Name, Artist, and/or LinkURL of the song.
// Parameters are the ID of the song and a Song object which contains the
// desired changes. Fields other than Name, Artist, and LinkURL are ignored.
func (s *SQLiteStore) UpdateSong(id int64, song Song) error {
	newName := song.Name
	newArtist := song.Artist
	newLinkURL := song.LinkURL

	_, err := s.db.ExecContext(s.ctx,
		"UPDATE songs SET name = $1, artist = $2, link_url = $3 WHERE id = $4",
		newName, newArtist, newLinkURL, id,
	)
	if err != nil {
		return fmt.Errorf("error updating song %d: %v", id, err)
	}

	return nil
}

// AddVote increments the Vote count on a Song with the given ID.
func (s *SQLiteStore) AddVote(id int64) error {
	song, err := s.GetSong(id)
	if err != nil {
		return err
	}

	votes := song.Votes + 1

	_, err = s.db.ExecContext(s.ctx,
		"UPDATE songs SET votes = $1 WHERE id = $2",
		votes, id,
	)
	if err != nil {
		return fmt.Errorf("error updating song %d: %v", id, err)
	}

	return nil
}

// Veto sets the Vetoed field of the Song with the given ID to true.
func (s *SQLiteStore) Veto(id int64) error {
	_, err := s.db.ExecContext(s.ctx,
		"UPDATE songs SET vetoed = $1 WHERE id = $2",
		true, id,
	)
	if err != nil {
		return fmt.Errorf("error updating song %d: %v", id, err)
	}

	return nil
}

// createTables creates the database tables if they do not already exist.
func (s *SQLiteStore) createTables() error {
	_, err := s.db.ExecContext(s.ctx,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			active BOOLEAN,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			vetoes INTEGER
		);
		CREATE TABLE IF NOT EXISTS songs (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			artist TEXT NOT NULL,
			link_url TEXT,
			votes INTEGER,
			vetoed BOOLEAN
		);
		CREATE TABLE IF NOT EXISTS votes (
			id INTEGER PRIMARY KEY,
			song INTEGER,
			voted_by INTEGER,
			FOREIGN KEY(song) REFERENCES songs(id),
			FOREIGN KEY(voted_by) REFERENCES users(id)
		)`,
	)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	return nil
}

// songExists queries the database for any song with the same name and artist
// as the given song and returns true if there is a match.
func (s *SQLiteStore) songExists(song Song) bool {
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

// userExists queries the database for any user with the same name as the
// given user and returns true if there is a match.
func (s *SQLiteStore) userExists(user User) bool {
	var name string

	err := s.db.QueryRowContext(s.ctx,
		"SELECT name FROM users WHERE name = $1",
		user.Name,
	).Scan(&name)

	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Fatalf("error checking for user %q: %v\n", user.Name, err)
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
