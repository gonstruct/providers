package encryption_test

import (
	"testing"

	"github.com/gonstruct/providers/encryption"
)

func TestEncrypt_String(t *testing.T) {
	fake := encryption.Fake()

	plain := "hello world"

	encrypted, err := encryption.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if encrypted == "" {
		t.Error("Encrypt() returned empty string")
	}

	if encrypted == plain {
		t.Error("Encrypt() returned plaintext")
	}

	fake.AssertEncryptedString(t, plain)
}

func TestEncrypt_Bytes(t *testing.T) {
	fake := encryption.Fake()

	plain := []byte("hello world")

	encrypted, err := encryption.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if encrypted == "" {
		t.Error("Encrypt() returned empty string")
	}

	fake.AssertEncryptedCount(t, 1)
}

func TestDecrypt_String(t *testing.T) {
	fake := encryption.Fake()

	plain := "secret message"

	encrypted, err := encryption.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	decrypted, err := encryption.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if string(decrypted) != plain {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plain)
	}

	fake.AssertDecryptedString(t, encrypted)
}

func TestDecryptString(t *testing.T) {
	fake := encryption.Fake()

	plain := "secret message"
	encrypted, _ := encryption.Encrypt(plain)

	decrypted, err := encryption.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString() error = %v", err)
	}

	if decrypted != plain {
		t.Errorf("DecryptString() = %q, want %q", decrypted, plain)
	}

	fake.AssertEncryptedCount(t, 1)
	fake.AssertDecryptedCount(t, 1)
}

func TestEncrypt_AssertNothingEncrypted(t *testing.T) {
	fake := encryption.Fake()

	// No operations
	fake.AssertNothingEncrypted(t)
	fake.AssertNothingDecrypted(t)
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	fake := encryption.Fake()

	plain := "round trip test"

	encrypted, err := encryption.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	decrypted, err := encryption.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if string(decrypted) != plain {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plain)
	}

	// Verify calls were tracked
	if fake.LastEncryptCall() == nil {
		t.Error("LastEncryptCall() should not be nil")
	}

	if fake.LastDecryptCall() == nil {
		t.Error("LastDecryptCall() should not be nil")
	}
}
