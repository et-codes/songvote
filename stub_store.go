package songvote

import "fmt"

// StubStore implements the SongStore interface, and keeps track of
// the calls against its methods. It is meant to be used for testing.
type StubStore struct {
	NextID            int64          // next song ID to be used
	GetSongCalls      []int64        // calls to to GetSong
	GetSongsCallCount int            // count of calls to GetSongs
	AddSongCalls      Songs          // calls to AddSong
	DeleteSongCalls   []int64        // calls to DeleteSong
	UpdateSongCalls   map[int64]Song // calls to UpdateSong
	AddVoteCalls      []int64        // calls to AddVote
	VetoCalls         []int64        // calls to Veto
}

// NewStubStore returns a reference to an empty StubStore.
func NewStubStore() *StubStore {
	return &StubStore{
		NextID:          1,
		UpdateSongCalls: make(map[int64]Song),
	}
}

func (s *StubStore) GetSong(id int64) (Song, error) {
	s.GetSongCalls = append(s.GetSongCalls, id)
	if id >= 10 {
		return Song{}, fmt.Errorf("song ID %d not found", id)
	}
	return Song{}, nil
}

func (s *StubStore) AddSong(song Song) (int64, error) {
	song.ID = s.NextID
	s.NextID++
	s.AddSongCalls = append(s.AddSongCalls, song)
	return song.ID, nil
}

func (s *StubStore) DeleteSong(id int64) error {
	s.DeleteSongCalls = append(s.DeleteSongCalls, id)
	if id >= 10 {
		return fmt.Errorf("song ID %d not found", id)
	}
	return nil
}

func (s *StubStore) GetSongs() Songs {
	s.GetSongsCallCount++
	return Songs{}
}

func (s *StubStore) UpdateSong(id int64, song Song) error {
	s.UpdateSongCalls[id] = song
	if id >= 10 {
		return fmt.Errorf("song ID %d not found", id)
	}
	return nil
}

func (s *StubStore) AddVote(id int64) error {
	s.AddVoteCalls = append(s.AddVoteCalls, id)
	if id >= 10 {
		return fmt.Errorf("song ID %d not found", id)
	}
	return nil
}

func (s *StubStore) Veto(id int64) error {
	s.VetoCalls = append(s.VetoCalls, id)
	if id >= 10 {
		return fmt.Errorf("song ID %d not found", id)
	}
	return nil
}