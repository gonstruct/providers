package aes_256_gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func (encrypter Adapter) Encrypt(plain []byte) (string, error) {
	block, err := aes.NewCipher(encrypter.Key())
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, plain, nil)
	final := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(final), nil
}
