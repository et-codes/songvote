package songvote

// Song contains information and methods about a single song.
type Song struct {
	ID       int64  `json:"id"`        // to be generated by store
	Name     string `json:"name"`      // song name
	Artist   string `json:"artist"`    // song artist
	LinkURL  string `json:"link_url"`  // link to YouTube, Spotify, etc.
	Votes    int    `json:"votes"`     // number of votes received
	VotedBy  Users  `json:"voted_by"`  // users that voted
	Vetoed   bool   `json:"vetoed"`    // this song has been vetoed
	VetoedBy User   `json:"vetoed_by"` // user that vetoed
}

// Songs contains a slice of Song objects.
type Songs []Song

// (Song).Equal returns whether names and artists of two Song objects are the same.
func (s Song) Equal(other Song) bool {
	return (s.Name == other.Name) && (s.Artist == other.Artist)
}
