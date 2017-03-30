package error

import "fmt"

// error <-- ServerError
// ServerError status/message
type ServerError interface {
	error
	Status() int32
	Message() string
}

type serverError struct {
	status  int32
	message string
}

// NewServerError create a new instance of ServerError
// with the specific status and message
func NewServerError(status int32, message string) ServerError {
	return &serverError{status, message}
}

// Status returns the status code of ServerError
func (e *serverError) Status() int32 {
	return e.status
}

// Message returns the message of ServerError
func (e *serverError) Message() string {
	return e.message
}

func (e *serverError) Error() string {
	return fmt.Sprintf("status: %d, message: %s", e.status, e.message)
}