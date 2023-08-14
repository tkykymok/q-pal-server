package errpkg

import "fmt"

type DatabaseError struct {
	InternalError error
	Operation     string
}

func (e *DatabaseError) Error() string {
	if e.InternalError != nil {
		return fmt.Sprintf("database error during %s: %v", e.Operation, e.InternalError)
	}
	return fmt.Sprintf("database error during %s", e.Operation)
}
