package common

import "fmt"

type ClientResponseError struct {
	HTTPCode int
	Message  string
}

func (e *ClientResponseError) Error() string {
	return fmt.Sprintf("[%d]: %s", e.HTTPCode, e.Message)
}
