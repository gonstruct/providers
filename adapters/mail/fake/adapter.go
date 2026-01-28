package fake

import (
	"context"
	"sync"

	"github.com/gonstruct/providers/entities"
)

// Adapter is a fake mail adapter for testing.
type Adapter struct {
	mu sync.RWMutex

	// Call tracking
	Calls []SendCall

	// Error injection
	SendError error

	// Custom send function
	SendFunc func(ctx context.Context, input entities.MailInput) error
}

type SendCall struct {
	To          []string
	From        string
	Subject     string
	HTML        string
	Attachments int
	Input       entities.MailInput
}

// New creates a new fake mail adapter.
func New() *Adapter {
	return &Adapter{}
}

// Reset clears all recorded calls.
func (a *Adapter) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Calls = nil
	a.SendError = nil
	a.SendFunc = nil
}

// --- Helper Methods ---

// SentCount returns the number of emails sent.
func (a *Adapter) SentCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.Calls)
}

// LastCall returns the last send call, or nil if none.
func (a *Adapter) LastCall() *SendCall {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.Calls) == 0 {
		return nil
	}

	return &a.Calls[len(a.Calls)-1]
}
