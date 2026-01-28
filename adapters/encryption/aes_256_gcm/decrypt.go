package aes_256_gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/gonstruct/providers/encryption"
)

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM
// Optional additionalData must match the AAD used during encryption
func (encrypter Adapter) Decrypt(base64Cipher string, additionalData ...[]byte) ([]byte, error) {
	// Use first additional data if provided, otherwise nil
	var aad []byte
	if len(additionalData) > 0 {
		aad = additionalData[0]
	}

	for _, key := range encrypter.Keys() {
		plain, err := encrypter.decrypt(key, base64Cipher, aad)
		if err == nil {
			return plain, nil
		}
	}

	return nil, encryption.Err("decrypt", encryption.ErrDecryptionFailed)
}

func (encrypter Adapter) decrypt(key []byte, base64Cipher string, additionalData []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64Cipher)
	if err != nil {
		return nil, encryption.Err("decode base64", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, encryption.Err("create cipher", encryption.ErrInvalidKey)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, encryption.Err("create GCM", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, encryption.Err("validate ciphertext", encryption.ErrCiphertextShort)
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, additionalData)
	if err != nil {
		return nil, encryption.Err("decrypt", err)
	}

	return plain, nil
}
