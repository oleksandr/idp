package entities

import "errors"

var (
	//ErrNotFound wraps not found error from lower data layer
	ErrNotFound = errors.New("Entity not found")
)
