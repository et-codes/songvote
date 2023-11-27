package songvote

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ServerError struct {
	Code int   `json:"code"`
	Err  error `json:"error"`
}

func NewServerError(code int, message string) ServerError {
	return ServerError{
		Code: code,
		Err:  fmt.Errorf(message),
	}
}

// Common errors
var (
	ErrMethod   = NewServerError(http.StatusMethodNotAllowed, "method not allowed")
	ErrConflict = NewServerError(http.StatusConflict, "resource already exists")
	ErrIDParse  = NewServerError(http.StatusBadRequest, "error parsing ID")
	ErrNotFound = NewServerError(http.StatusNotFound, "resource not found")
)

func (e ServerError) Error() string {
	return fmt.Sprintf("status: %d, error: %v", e.Code, e.Err)
}

func writeError(w http.ResponseWriter, e ServerError) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	json, _ := json.Marshal(e)
	fmt.Fprint(w, json)
}
