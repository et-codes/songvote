package songvote

// SQLSongStore is a song store backed by a SQL database.
type SQLSongStore struct{}

// NewSQLSongStore returns a pointer to a newly initialized store.
func NewSQLSongStore() *SQLSongStore {
	return &SQLSongStore{}
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
