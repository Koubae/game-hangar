package common

import "fmt"

type BusinessError struct {
	HTTPCode int
	Message  string
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("[%d]: %s", e.HTTPCode, e.Message)
}
