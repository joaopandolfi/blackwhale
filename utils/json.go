package utils

import (
	"encoding/json"
	"fmt"
)

func JsonString(m interface{}) (string, error) {
	if m == nil {
		return "", fmt.Errorf("m can not be null")
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshalling data: %w", err)
	}

	return string(b), nil
}
