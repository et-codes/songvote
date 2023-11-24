package songvote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Song contains information about a single song.
type Song struct {
	ID      int64  `json:"id"`       // to be generated by store
	Name    string `json:"name"`     // song name
	Artist  string `json:"artist"`   // song artist
	LinkURL string `json:"link_url"` // link to YouTube, Spotify, etc.
	Votes   int    `json:"votes"`    // number of votes received
	Vetoed  bool   `json:"vetoed"`   // this song has been vetoed
}

// (Song).Equal returns whether names and artists of two Song objects are the same.
func (s Song) Equal(other Song) bool {
	return (s.Name == other.Name) && (s.Artist == other.Artist)
}

// MarshalSong returns JSON-encoded string of the Song object.
func MarshalSong(song Song) (string, error) {
	output := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(output).Encode(song)
	if err != nil {
		return "", fmt.Errorf("problem encoding song to JSON: %w", err)
	}
	return output.String(), nil
}

// UnmarshalSong returns JSON-encoded string of the Song object.
func UnmarshalSong(input io.Reader, song *Song) error {
	err := json.NewDecoder(input).Decode(song)
	if err != nil {
		return fmt.Errorf("problem decoding song from JSON: %w", err)
	}
	return nil
}
