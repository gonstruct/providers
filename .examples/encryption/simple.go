package main

import (
	"github.com/gonstruct/providers/adapters/encryption/aes_256_gcm"
	"github.com/gonstruct/providers/encryption"
)

func main() {
	key, err := encryption.GenerateKey(encryption.KeyFormatBase64, encryption.WithAdapter(&aes_256_gcm.Adapter{}))
	if err != nil {
		panic(err)
	}

	println("Generated Key:", key)

	encryption.Adapt(
		&aes_256_gcm.Adapter{
			Key: func() []byte {
				return encryption.ParseKeys(key)[0]
			},
			PreviousKeys: func() [][]byte {
				return [][]byte{
					[]byte("old example key 123old example key 123"),
				}
			},
		},
	)

	base64cipher, err := encryption.Encrypt("hello world")
	if err != nil {
		panic(err)
	}

	println("Encrypted:", base64cipher)

	plaintext, err := encryption.DecryptString(base64cipher)
	if err != nil {
		panic(err)
	}

	println("Decrypted:", plaintext)
}
