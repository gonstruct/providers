package contracts

type Encryption interface {
	// Encrypt encrypts plain bytes, optionally with additional authenticated data (AAD)
	Encrypt(plain []byte, additionalData ...[]byte) (string, error)
	// Decrypt decrypts base64-encoded ciphertext, optionally verifying AAD
	Decrypt(base64Cipher string, additionalData ...[]byte) ([]byte, error)
	// GenerateKey generates a new encryption key
	GenerateKey() ([]byte, error)
}
