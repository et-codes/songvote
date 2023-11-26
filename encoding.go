package songvote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// MarshalJSON returns JSON-encoded string of an object and an error if there is one.
func MarshalJSON(obj any) (string, error) {
	output := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(output).Encode(obj)
	if err != nil {
		return "", fmt.Errorf("problem encoding to JSON: %w", err)
	}
	return output.String(), nil
}

// UnmarshalJSON decodes the JSON string into the references struct and returns
// an error if there is one.
func UnmarshalJSON[T any](input io.Reader, obj *T) error {
	err := json.NewDecoder(input).Decode(obj)
	if err != nil {
		return fmt.Errorf("problem decoding from JSON: %w", err)
	}
	return nil
}
