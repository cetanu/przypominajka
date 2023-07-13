package models

import "errors"

var ErrNotFound = errors.New("not found")

// Validation errors
var (
	ErrMissingNameOrNames = errors.New("'name' or 'names' must be provided")
	ErrNameOrNames        = errors.New("'name' is mutually exclusive with 'names'")
	ErrNamesArePair       = errors.New("'names' must have two elements")
	ErrInvalidEventType   = errors.New("invalid event type")
)
