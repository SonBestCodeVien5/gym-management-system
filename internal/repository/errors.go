package repository

import "errors"

// ErrNotFound is a storage-agnostic not-found error exposed to upper layers.
var ErrNotFound = errors.New("resource not found")
