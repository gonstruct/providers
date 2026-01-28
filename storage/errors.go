package storage

import (
	"errors"
	"fmt"
)

// Sentinel errors for storage operations.
var (
	ErrFileNotFound      = errors.New("file not found")
	ErrDirectoryNotFound = errors.New("directory not found")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrInvalidPath       = errors.New("invalid path")
	ErrAlreadyExists     = errors.New("file already exists")
)

// Err wraps an error with storage context.
func Err(op string, err error) error {
	return fmt.Errorf("storage: %s: %w", op, err)
}

// PathErr wraps an error with storage context and path info.
func PathErr(op, path string, err error) error {
	return fmt.Errorf("storage: %s %q: %w", op, path, err)
}
