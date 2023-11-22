package songvote

import (
	"fmt"
	"net/http"
)

func SongVoteServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Would?")
}
