package store

import "errors"

var (
	// ErrRecordNotFound error tells us, that we can't fing target record in DB
	ErrRecordNotFound = errors.New("record not found")
)
