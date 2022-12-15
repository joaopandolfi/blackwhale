package handlers

import (
	"net/http"

	"github.com/go-playground/validator"
)

// defaultValidator its a struct validator
var defaultValidator = validator.New()

// UnmarshalSnakeCaseAndValidate -
// Unmarshal payload using snake case and validade using default validator
func UnmarshalSnakeCaseAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) (string, error) {

	err := SnakeCaseDecoder(r.Body).Decode(v)
	if err != nil {
		return "unmarshaling", err
	}

	err = defaultValidator.Struct(v)
	if err != nil {
		return "validating body", err
	}

	return "", nil
}
