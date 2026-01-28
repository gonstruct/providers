package fake

import (
	"testing"
)

// AssertEncrypted asserts that data was encrypted at least once
func (a *Adapter) AssertEncrypted(t testing.TB) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.EncryptCalls) == 0 {
		t.Error("Expected data to be encrypted, but it was not")
	}
}

// AssertEncryptedCount asserts the exact number of encryptions
func (a *Adapter) AssertEncryptedCount(t testing.TB, count int) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.EncryptCalls) != count {
		t.Errorf("Expected %d encryptions, got %d", count, len(a.EncryptCalls))
	}
}

// AssertEncryptedString asserts that a specific string was encrypted
func (a *Adapter) AssertEncryptedString(t testing.TB, plaintext string) {
	t.Helper()
	a.AssertEncryptedData(t, []byte(plaintext))
}

// AssertEncryptedData asserts that specific data was encrypted
func (a *Adapter) AssertEncryptedData(t testing.TB, plaintext []byte) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, call := range a.EncryptCalls {
		if bytesEqual(call.Plaintext, plaintext) {
			return
		}
	}

	t.Errorf("Expected data %q to be encrypted, but it was not", plaintext)
}

// AssertDecrypted asserts that data was decrypted at least once
func (a *Adapter) AssertDecrypted(t testing.TB) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.DecryptCalls) == 0 {
		t.Error("Expected data to be decrypted, but it was not")
	}
}

// AssertDecryptedCount asserts the exact number of decryptions
func (a *Adapter) AssertDecryptedCount(t testing.TB, count int) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.DecryptCalls) != count {
		t.Errorf("Expected %d decryptions, got %d", count, len(a.DecryptCalls))
	}
}

// AssertDecryptedString asserts that a specific ciphertext was decrypted
func (a *Adapter) AssertDecryptedString(t testing.TB, ciphertext string) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, call := range a.DecryptCalls {
		if call.Ciphertext == ciphertext {
			return
		}
	}

	t.Errorf("Expected ciphertext %q to be decrypted, but it was not", ciphertext)
}

// AssertNothingEncrypted asserts that no encryptions occurred
func (a *Adapter) AssertNothingEncrypted(t testing.TB) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.EncryptCalls) > 0 {
		t.Errorf("Expected no encryptions, but %d occurred", len(a.EncryptCalls))
	}
}

// AssertNothingDecrypted asserts that no decryptions occurred
func (a *Adapter) AssertNothingDecrypted(t testing.TB) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.DecryptCalls) > 0 {
		t.Errorf("Expected no decryptions, but %d occurred", len(a.DecryptCalls))
	}
}
