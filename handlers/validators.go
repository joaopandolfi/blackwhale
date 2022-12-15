package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/joaopandolfi/blackwhale/handlers/errors"
	"github.com/joaopandolfi/blackwhale/utils"
)

// defaultValidator its a struct validator
var defaultValidator = validator.New()

// UnmarshalSnakeCaseAndValidate -
// Unmarshal payload using snake case and validade using default validator
func UnmarshalSnakeCaseAndValidate(w http.ResponseWriter, r *http.Request, functionName string, v interface{}) error {

	err := SnakeCaseDecoder(r.Body).Decode(v)
	if err != nil {
		utils.Error(fmt.Sprintf("%s - unmarshaling", functionName), err.Error())
		ResponseTypedError(w, errors.ErrorCodeInvalidBody, errors.ErrorMessageInvalidBody, err)
		return err
	}

	err = defaultValidator.Struct(v)
	if err != nil {
		utils.Error(fmt.Sprintf("%s - validating body", functionName), err.Error())
		ResponseTypedError(w, errors.ErrorCodeInvalidBody, errors.ErrorMessageInvalidBody, err)
		return err
	}

	return nil
}
