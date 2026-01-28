package encryption

import (
	"github.com/gonstruct/providers/adapters/encryption/aes_256_gcm"
	"github.com/gonstruct/providers/encryption"
)

func Example() {
	// Set up encryption with AES-256-GCM
	encryption.Adapt(&aes_256_gcm.Adapter{
		Key: func() []byte {
			return []byte("your-32-byte-secret-key-here!!")
		},
	})

	// Encrypt data
	encrypted, _ := encryption.Encrypt("sensitive data")

	// Decrypt data
	plaintext, _ := encryption.DecryptString(encrypted)

	println(plaintext) // "sensitive data"
}
