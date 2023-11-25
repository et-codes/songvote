package songvote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) ToJSON() string {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(e); err != nil {
		log.Printf("Problem converting error to JSON: %v\n", err)
		return fmt.Sprintf(`{"code":"%d","message":"%s"}`, e.Code, e.Message)
	}
	return buf.String()
}
