package contracts

type Encryption interface {
	Encrypt(plain []byte) (string, error)
	Decrypt(base64Cipher string) ([]byte, error)
	GenerateKey() ([]byte, error)
}
