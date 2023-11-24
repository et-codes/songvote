package songvote

import (
	"fmt"
)

type InMemorySongStore struct {
	songs  []Song // slice of songs in the store
	nextID int64  // next serialized song ID to use
}

func NewInMemorySongStore() *InMemorySongStore {
	return &InMemorySongStore{
		songs:  []Song{},
		nextID: 1,
	}
}

func (i *InMemorySongStore) GetSong(id int64) (Song, error) {
	for _, song := range i.songs {
		if song.ID == id {
			return song, nil
		}
	}
	return Song{}, fmt.Errorf("song ID %d not found", id)
}

func (i *InMemorySongStore) GetSongs() []Song {
	return i.songs
}

func (i *InMemorySongStore) AddSong(song Song) (int64, error) {
	if i.songExists(song) {
		return 0, fmt.Errorf("%q by %q already in store", song.Name, song.Artist)
	}
	song.ID = i.nextID
	i.nextID++
	i.songs = append(i.songs, song)
	return song.ID, nil
}

func (i *InMemorySongStore) DeleteSong(id int64) error {
	return fmt.Errorf("not implemented")
}

func (i *InMemorySongStore) songExists(song Song) bool {
	for _, s := range i.songs {
		if song.Equal(s) {
			return true
		}
	}
	return false
}
