package songvote

type InMemorySongStore struct{}

func (i *InMemorySongStore) GetSong(name string) string {
	return "What What"
}

func (i *InMemorySongStore) AddSong(name string) {}
