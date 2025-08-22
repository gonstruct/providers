package aes_256_gcm

import (
	"crypto/rand"
)

const (
	KeySize = 32 // 256 bits
)

func (encrypter Adapter) GenerateKey() ([]byte, error) {
	key := make([]byte, KeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	return key, nil
}
