package database

import "errors"

type NotFoundError error

func NewNotFoundError() NotFoundError {
	err := errors.New("")
	return NotFoundError(err)
}
