package encryption

import (
	"github.com/gonstruct/providers/adapters/encryption/fake"
)

// Fake sets up a fake encryption adapter for testing and returns it for assertions.
// This replaces any existing encryption provider.
//
// Example:
//
//	func TestEncryptData(t *testing.T) {
//	    fake := encryption.Fake()
//
//	    // Your code that uses encryption.Encrypt(), encryption.Decrypt()
//	    encrypted, _ := encryption.Encrypt(ctx, []byte("secret"))
//
//	    // Assert
//	    fake.AssertEncrypted(t)
//	    fake.AssertEncryptedData(t, []byte("secret"))
//	}
func Fake() *fake.Adapter {
	adapter := fake.New()

	globalProvider = &provider{
		adapter: adapter,
	}

	return adapter
}
