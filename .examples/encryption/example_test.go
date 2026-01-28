package encryption

import (
	"testing"

	"github.com/gonstruct/providers/encryption"
)

func TestEncryptDecrypt(t *testing.T) {
	fake := encryption.Fake()

	encrypted, _ := encryption.Encrypt("secret")
	decrypted, _ := encryption.DecryptString(encrypted)

	if decrypted != "secret" {
		t.Errorf("got %q, want %q", decrypted, "secret")
	}

	fake.AssertEncrypted(t, "secret")
	fake.AssertDecrypted(t, encrypted)
}

func TestNothingEncrypted(t *testing.T) {
	fake := encryption.Fake()

	fake.AssertNothingEncrypted(t)
}
