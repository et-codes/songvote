package songvote

import (
	"fmt"
)

type InMemorySongStore struct {
	songs  []Song // slice of songs in the store
	nextId int    // next serialized song ID to use
}

func NewInMemorySongStore() *InMemorySongStore {
	return &InMemorySongStore{
		songs:  []Song{},
		nextId: 1,
	}
}

func (i *InMemorySongStore) GetSong(id int) (Song, error) {
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

func (i *InMemorySongStore) AddSong(song Song) (int, error) {
	if i.songExists(song) {
		return 0, fmt.Errorf("%q by %q already in store", song.Name, song.Artist)
	}
	song.ID = i.nextId
	i.nextId++
	i.songs = append(i.songs, song)
	return song.ID, nil
}

func (i *InMemorySongStore) songExists(song Song) bool {
	for _, s := range i.songs {
		if song.Equal(s) {
			return true
		}
	}
	return false
}
