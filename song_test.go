package songvote_test

import (
	"testing"

	"github.com/et-codes/songvote"
)

func TestSongEquals(t *testing.T) {
	song1 := songvote.Song{Name: "Snot", Artist: "Snot"}
	song2 := songvote.Song{Name: "Snot", Artist: "Snot"}
	assertTrue(t, song1.Equals(song2))
}

func TestSongMarshal(t *testing.T) {
	song := songvote.Song{
		Name:    "Alive",
		Artist:  "Pearl Jam",
		LinkURL: `https://youtu.be/qM0zINtulhM`,
		Votes:   3,
		Vetoed:  false,
	}

	want := `{"id":0,"name":"Alive","artist":"Pearl Jam","link_url":"https://youtu.be/qM0zINtulhM","votes":3,"vetoed":false}` + "\n"

	got, err := song.Marshal()
	assertNoError(t, err)
	assertEqual(t, got, want)
}

func TestSongUnmarshal(t *testing.T) {
	jsonString := `{"id":0,"name":"Alive","artist":"Pearl Jam","link_url":"https://youtu.be/qM0zINtulhM","votes":3,"vetoed":false}` + "\n"

	want := songvote.Song{
		Name:    "Alive",
		Artist:  "Pearl Jam",
		LinkURL: `https://youtu.be/qM0zINtulhM`,
		Votes:   3,
		Vetoed:  false,
	}
	got := songvote.Song{}

	err := songvote.UnmarshalSong(jsonString, &got)
	assertNoError(t, err)
	assertEqual(t, got, want)
}
