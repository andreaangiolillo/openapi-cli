package errors

import "fmt"

type SchemaConflictError struct {
	Entry string
}

func (e SchemaConflictError) Error() string {
	return fmt.Sprintf("there was a conflict on a the Schema component: %s", e.Entry)
}
