package fake

import (
	"encoding/base64"
)

func (a *Adapter) Encrypt(plain []byte, additionalData ...[]byte) (string, error) {
	if a.EncryptFunc != nil {
		result, err := a.EncryptFunc(plain, additionalData...)
		if err == nil {
			a.mu.Lock()
			a.EncryptCalls = append(a.EncryptCalls, EncryptCall{
				Plaintext:      plain,
				AdditionalData: additionalData,
				Result:         result,
			})
			a.mu.Unlock()
		}

		return result, err
	}

	if a.EncryptError != nil {
		return "", a.EncryptError
	}

	// Simple fake encryption: base64 encode with prefix
	result := "fake:" + base64.StdEncoding.EncodeToString(plain)

	a.mu.Lock()
	a.EncryptCalls = append(a.EncryptCalls, EncryptCall{Plaintext: plain, AdditionalData: additionalData, Result: result})
	a.store[result] = plain
	a.mu.Unlock()

	return result, nil
}
