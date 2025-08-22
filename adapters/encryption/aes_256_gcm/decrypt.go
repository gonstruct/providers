package aes_256_gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func (encrypter Adapter) Decrypt(base64Cipher string) ([]byte, error) {
	for _, key := range encrypter.Keys() {
		plain, err := encrypter.decrypt(key, base64Cipher)
		if err == nil {
			return plain, nil
		}
	}

	return nil, errors.New("decryption failed with all provided keys")
}

func (encrypter Adapter) decrypt(key []byte, base64Cipher string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64Cipher)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}
