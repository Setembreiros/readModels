package database

import "fmt"

type NotFoundError struct {
	table string
	key   any
}

func (e *NotFoundError) Error() string {
	errorMessage := fmt.Sprintf("Data in table %s not found for key %v", e.table, e.key)
	return errorMessage
}

func NewNotFoundError(table string, key any) *NotFoundError {
	return &NotFoundError{
		table: table,
		key:   key,
	}
}

type InvalidResultsError struct {
	expected string
	got      string
}

func (e *InvalidResultsError) Error() string {
	return fmt.Sprintf("Invalid results parameter: expected %s, got %s", e.expected, e.got)
}

func NewInvalidResultsError(expected, got string) *InvalidResultsError {
	return &InvalidResultsError{
		expected: expected,
		got:      got,
	}
}

func NewInvalidSlicePointerError(gotType string) *InvalidResultsError {
	return NewInvalidResultsError("pointer to slice", gotType)
}
