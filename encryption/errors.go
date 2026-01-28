package encryption

import (
	"errors"
	"fmt"
)

// Sentinel errors for encryption operations
var (
	ErrDecryptionFailed = errors.New("decryption failed")
	ErrInvalidKey       = errors.New("invalid encryption key")
	ErrCiphertextShort  = errors.New("ciphertext too short")
)

// Err wraps an error with encryption context
func Err(op string, err error) error {
	return fmt.Errorf("encryption: %s: %w", op, err)
}
