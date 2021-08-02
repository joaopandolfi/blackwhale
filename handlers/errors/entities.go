package errors

// TypedError is used to send to interface errors with code
type TypedError struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Stack   interface{} `json:"stack,omitempty"`
	Success bool        `json:"success"`
}

// NewTypedError constructor
func NewTypedError(code int, message string, stack error) *TypedError {
	var st string
	if stack != nil {
		st = stack.Error()
	}
	return &TypedError{
		Code:    code,
		Message: message,
		Stack:   st,
		Success: false,
	}
}

//AppError -
type AppError struct { //nolint
	HTTPCode int    `json:"-"`
	Code     string `json:"errCode"`
	Message  string `json:"errMessage,omitempty"`
}

func (e AppError) Error() string {
	return e.Message
}
