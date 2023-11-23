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
	_, err := i.findSongByID(song.ID)
	if err == nil {
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
