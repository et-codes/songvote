package songvote

import (
	"fmt"
	"net/http"
)

func Server(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Would?")
}
