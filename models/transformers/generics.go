package transformers

import (
	"encoding/json"
	"fmt"
)

// ConvertGeneric - converts interface inside expected struct using json serialization
func ConvertTo[T any](v interface{}) (T, error) {
	var result T

	b, err := json.Marshal(v)
	if err != nil {
		return result, fmt.Errorf("marshaling interface: %w", err)
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return result, fmt.Errorf("unmarshaling interface: %w", err)
	}

	return result, nil
}

// In returns the index of the first occurrence of v in s, or -1 if not present.
func In[E comparable](s []E, v E) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}
