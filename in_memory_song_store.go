package songvote

import "fmt"

type InMemorySongStore struct {
	songs []Song
}

func (i *InMemorySongStore) GetSong(id int) Song {
	song, err := i.findSongByID(id)
	if err != nil {
		return Song{}
	}
	return song
}

func (i *InMemorySongStore) AddSong(song Song) {
	if !i.songExists(song) {
		i.songs = append(i.songs, song)
	}
}

func (i *InMemorySongStore) findSongByID(id int) (Song, error) {
	for _, song := range i.songs {
		if song.ID == id {
			return song, nil
		}
	}
	return Song{}, fmt.Errorf("song ID %d not found", id)
}

func (i *InMemorySongStore) songExists(song Song) bool {
	for _, s := range i.songs {
		if song.Equals(s) {
			return true
		}
	}
	return false
}
