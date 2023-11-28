package songvote

import "fmt"

// StubStore implements the SongStore interface, and keeps track of
// the calls against its methods. It is meant to be used for testing.
type StubStore struct {
	NextSongID        int64           // next song ID to be used
	AddSongCalls      Songs           // calls to AddSong
	GetSongsCallCount int             // count of calls to GetSongs
	GetSongCalls      []int64         // calls to to GetSong
	DeleteSongCalls   []int64         // calls to DeleteSong
	UpdateSongCalls   map[int64]Song  // calls to UpdateSong
	AddVoteCalls      Votes           // calls to AddVote
	VetoCalls         map[int64]int64 // calls to Veto, [songID]userID

	NextUserID        int64          // next user ID to be used
	AddUserCalls      Users          // calls to AddUser
	GetUsersCallCount int            // count of calls to GetUsers
	GetUserCalls      []int64        // calls to to GetUser
	DeleteUserCalls   []int64        // calls to DeleteUser
	UpdateUserCalls   map[int64]User // calls to UpdateUser
}

// NewStubStore returns a reference to an empty StubStore.
func NewStubStore() *StubStore {
	return &StubStore{
		NextSongID:      1,
		UpdateSongCalls: make(map[int64]Song),
		VetoCalls:       make(map[int64]int64),
		NextUserID:      1,
		UpdateUserCalls: make(map[int64]User),
	}
}

func (s *StubStore) AddUser(user User) (int64, error) {
	user.ID = s.NextUserID
	s.NextUserID++
	s.AddUserCalls = append(s.AddUserCalls, user)
	return user.ID, nil
}

func (s *StubStore) GetUsers() Users {
	s.GetUsersCallCount++
	return Users{}
}

func (s *StubStore) GetUser(id int64) (User, error) {
	s.GetUserCalls = append(s.GetUserCalls, id)
	if id >= 10 {
		return User{}, fmt.Errorf("user ID %d not found", id)
	}
	return User{}, nil
}

func (s *StubStore) DeleteUser(id int64) error {
	s.DeleteUserCalls = append(s.DeleteUserCalls, id)
	if id >= 10 {
		return fmt.Errorf("user ID %d not found", id)
	}
	return nil
}

func (s *StubStore) UpdateUser(id int64, user User) error {
	s.UpdateUserCalls[id] = user
	if id >= 10 {
		return fmt.Errorf("user ID %d not found", id)
	}
	return nil
}

func (s *StubStore) GetSong(id int64) (Song, error) {
	s.GetSongCalls = append(s.GetSongCalls, id)
	if id >= 10 {
		return Song{}, fmt.Errorf("song ID %d not found", id)
	}
	return Song{}, nil
}

func (s *StubStore) AddSong(song Song) (int64, error) {
	song.ID = s.NextSongID
	s.NextSongID++
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

func (s *StubStore) AddVote(vote Vote) error {
	s.AddVoteCalls = append(s.AddVoteCalls, vote)
	if vote.SongID >= 10 {
		return fmt.Errorf("song ID %d not found", vote.SongID)
	}
	return nil
}

func (s *StubStore) Veto(songID, userID int64) error {
	s.VetoCalls[songID] = userID
	if songID >= 10 {
		return fmt.Errorf("song ID %d not found", songID)
	}
	return nil
}
