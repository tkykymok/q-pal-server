package errors

import "fmt"

type UnexpectedError struct {
	InternalError error
	Operation     string
}

func (e *UnexpectedError) Error() string {
	if e.InternalError != nil {
		return fmt.Sprintf("an unexpected error occurred during %s: %v", e.Operation, e.InternalError)
	}
	return fmt.Sprintf("an unexpected error occurred during %s", e.Operation)
}
