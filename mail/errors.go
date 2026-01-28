package mail

import (
	"errors"
	"fmt"
)

// Sentinel errors for mail operations
var (
	ErrNoSubject    = errors.New("no subject specified")
	ErrNoSender     = errors.New("no sender specified")
	ErrNoRecipients = errors.New("no recipients specified")
	ErrSendFailed   = errors.New("failed to send email")
)

// Err wraps an error with mail context
func Err(op string, err error) error {
	return fmt.Errorf("mail: %s: %w", op, err)
}
