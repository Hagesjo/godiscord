package discordgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrInvalidJSON = errors.New("invalid json")

// ReadJSON is a convenience function to unmarshal a client read into a json struct.
// Ideally, this should be a method to the client, but due to go limitations with generic methods, we can't do that for now.
func UnmarshalJSON[T any](data json.RawMessage) (T, error) {
	var ret T
	if err := json.Unmarshal(data, &ret); err != nil {
		return ret, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return ret, nil
}

func MustUnmarshalJSON[T any](data json.RawMessage) T {
	ret, err := UnmarshalJSON[T](data)
	if err != nil {
		panic(err)
	}

	return ret
}
