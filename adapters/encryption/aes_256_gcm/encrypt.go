package aes_256_gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/gonstruct/providers/encryption"
)

// Encrypt encrypts plain bytes using AES-256-GCM
// Optional additionalData provides authenticated data (AAD) that is verified but not encrypted
func (encrypter Adapter) Encrypt(plain []byte, additionalData ...[]byte) (string, error) {
	block, err := aes.NewCipher(encrypter.Key())
	if err != nil {
		return "", encryption.Err("create cipher", encryption.ErrInvalidKey)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", encryption.Err("create GCM", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", encryption.Err("generate nonce", err)
	}

	// Use first additional data if provided, otherwise nil
	var aad []byte
	if len(additionalData) > 0 {
		aad = additionalData[0]
	}

	ciphertext := gcm.Seal(nil, nonce, plain, aad)
	final := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(final), nil
}
