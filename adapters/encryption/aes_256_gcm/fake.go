package aes_256_gcm

// FakeAdapter is a mock encryption adapter for testing.
type FakeAdapter struct {
	// EncryptFunc allows customizing the Encrypt behavior
	EncryptFunc func(plain []byte, additionalData ...[]byte) (string, error)
	// DecryptFunc allows customizing the Decrypt behavior
	DecryptFunc func(base64Cipher string, additionalData ...[]byte) ([]byte, error)
	// GenerateKeyFunc allows customizing the GenerateKey behavior
	GenerateKeyFunc func() ([]byte, error)

	// EncryptCalls records all calls to Encrypt
	EncryptCalls []FakeEncryptCall
	// DecryptCalls records all calls to Decrypt
	DecryptCalls []FakeDecryptCall
	// GenerateKeyCalls records the number of calls to GenerateKey
	GenerateKeyCalls int
}

type FakeEncryptCall struct {
	Plain          []byte
	AdditionalData [][]byte
}

type FakeDecryptCall struct {
	Base64Cipher   string
	AdditionalData [][]byte
}

// Fake creates a new mock encryption adapter with default behaviors.
func Fake() *FakeAdapter {
	return &FakeAdapter{
		EncryptFunc: func(plain []byte, additionalData ...[]byte) (string, error) {
			return "fake-encrypted-" + string(plain), nil
		},
		DecryptFunc: func(base64Cipher string, additionalData ...[]byte) ([]byte, error) {
			return []byte(base64Cipher), nil
		},
		GenerateKeyFunc: func() ([]byte, error) {
			return []byte("fake-key-32-bytes-for-testing!!!"), nil
		},
	}
}

func (a *FakeAdapter) Encrypt(plain []byte, additionalData ...[]byte) (string, error) {
	a.EncryptCalls = append(a.EncryptCalls, FakeEncryptCall{
		Plain:          plain,
		AdditionalData: additionalData,
	})

	return a.EncryptFunc(plain, additionalData...)
}

func (a *FakeAdapter) Decrypt(base64Cipher string, additionalData ...[]byte) ([]byte, error) {
	a.DecryptCalls = append(a.DecryptCalls, FakeDecryptCall{
		Base64Cipher:   base64Cipher,
		AdditionalData: additionalData,
	})

	return a.DecryptFunc(base64Cipher, additionalData...)
}

func (a *FakeAdapter) GenerateKey() ([]byte, error) {
	a.GenerateKeyCalls++

	return a.GenerateKeyFunc()
}

// Reset clears all recorded calls.
func (a *FakeAdapter) Reset() {
	a.EncryptCalls = nil
	a.DecryptCalls = nil
	a.GenerateKeyCalls = 0
}
