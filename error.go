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
	// Method Not Allowed (405)
	ErrMethod = NewServerError(http.StatusMethodNotAllowed, "method not allowed")
	// Conflict (409) - resource already exists
	ErrConflict = NewServerError(http.StatusConflict, "resource already exists")
	// Bad Request (400) - error parsing ID
	ErrIDParse = NewServerError(http.StatusBadRequest, "error parsing ID")
	// Not Found (404)
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
