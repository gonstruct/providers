package fake

import (
	"encoding/base64"
)

func (a *Adapter) Decrypt(base64Cipher string, additionalData ...[]byte) ([]byte, error) {
	if a.DecryptFunc != nil {
		result, err := a.DecryptFunc(base64Cipher, additionalData...)
		if err == nil {
			a.mu.Lock()
			a.DecryptCalls = append(a.DecryptCalls, DecryptCall{Ciphertext: base64Cipher, AdditionalData: additionalData, Result: result})
			a.mu.Unlock()
		}
		return result, err
	}

	if a.DecryptError != nil {
		return nil, a.DecryptError
	}

	// Look up in store
	a.mu.RLock()
	plaintext, ok := a.store[base64Cipher]
	a.mu.RUnlock()

	if !ok {
		// Try stripping the prefix and decoding
		prefix := "fake:"
		if len(base64Cipher) > len(prefix) && base64Cipher[:len(prefix)] == prefix {
			decoded, err := base64.StdEncoding.DecodeString(base64Cipher[len(prefix):])
			if err == nil {
				plaintext = decoded
			}
		}
	}

	a.mu.Lock()
	a.DecryptCalls = append(a.DecryptCalls, DecryptCall{Ciphertext: base64Cipher, AdditionalData: additionalData, Result: plaintext})
	a.mu.Unlock()

	return plaintext, nil
}
