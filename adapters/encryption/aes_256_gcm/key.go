package aes_256_gcm

import (
	"crypto/rand"
)

func (encrypter Adapter) GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256 requires a 32-byte key
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}
