package songvote

import (
	"encoding/json"
	"fmt"
	"io"
)

// MarshalJSON returns JSON-encoded string of an object and an error if there is one.
func MarshalJSON(obj any) (string, error) {
	output, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("problem marshalling JSON: %w", err)
	}
	return string(output), nil
}

// UnmarshalJSON decodes the JSON string into the referenced struct and returns
// an error if there is one.
func UnmarshalJSON(in io.Reader, obj any) error {
	data, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("problem reading input: %w", err)
	}
	if err := json.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("problem unmarshalling JSON: %w", err)
	}
	return nil
}
