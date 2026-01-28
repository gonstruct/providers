package fake

import (
	"sync"
)

// Adapter is a fake encryption adapter for testing.
type Adapter struct {
	mu sync.RWMutex

	// Call tracking
	EncryptCalls []EncryptCall
	DecryptCalls []DecryptCall

	// Error injection
	EncryptError     error
	DecryptError     error
	GenerateKeyError error

	// Custom functions
	EncryptFunc     func(plain []byte, additionalData ...[]byte) (string, error)
	DecryptFunc     func(base64Cipher string, additionalData ...[]byte) ([]byte, error)
	GenerateKeyFunc func() ([]byte, error)

	// Simple store for encrypt/decrypt round-trip
	store map[string][]byte
}

type EncryptCall struct {
	Plaintext      []byte
	AdditionalData [][]byte
	Result         string
}

type DecryptCall struct {
	Ciphertext     string
	AdditionalData [][]byte
	Result         []byte
}

// New creates a new fake encryption adapter.
func New() *Adapter {
	return &Adapter{
		store: make(map[string][]byte),
	}
}

// Reset clears all recorded calls.
func (a *Adapter) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.EncryptCalls = nil
	a.DecryptCalls = nil
	a.EncryptError = nil
	a.DecryptError = nil
	a.GenerateKeyError = nil
	a.EncryptFunc = nil
	a.DecryptFunc = nil
	a.GenerateKeyFunc = nil
	a.store = make(map[string][]byte)
}

// --- Helper Methods ---

// EncryptCount returns the number of encryption calls.
func (a *Adapter) EncryptCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.EncryptCalls)
}

// DecryptCount returns the number of decryption calls.
func (a *Adapter) DecryptCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.DecryptCalls)
}

// LastEncryptCall returns the last encrypt call, or nil if none.
func (a *Adapter) LastEncryptCall() *EncryptCall {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.EncryptCalls) == 0 {
		return nil
	}

	return &a.EncryptCalls[len(a.EncryptCalls)-1]
}

// LastDecryptCall returns the last decrypt call, or nil if none.
func (a *Adapter) LastDecryptCall() *DecryptCall {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.DecryptCalls) == 0 {
		return nil
	}

	return &a.DecryptCalls[len(a.DecryptCalls)-1]
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
