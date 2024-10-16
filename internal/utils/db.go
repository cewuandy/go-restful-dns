package utils

import (
	"encoding/json"
	"fmt"
)

// Convert This function represents convert struct with JSON
func Convert(from interface{}, to interface{}) error {
	marshalled, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("JSON marshal error: %w", err)
	}

	err = json.Unmarshal(marshalled, to)
	if err != nil {
		return fmt.Errorf("JSON unmarshal error: %w", err)
	}

	return nil
}
