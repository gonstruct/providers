package aes_256_gcm_test

import (
	"testing"

	aes "github.com/gonstruct/providers/adapters/encryption/aes_256_gcm"
)

func TestGenerateKey(t *testing.T) {
	adapter := aes.Adapter{
		Key: func() []byte { return make([]byte, 32) },
	}

	key, err := adapter.GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}

	if len(key) != aes.KeySize {
		t.Errorf("GenerateKey() key length = %d, want %d", len(key), aes.KeySize)
	}
}

func TestGenerateKey_Uniqueness(t *testing.T) {
	adapter := aes.Adapter{
		Key: func() []byte { return make([]byte, 32) },
	}

	keys := make(map[string]bool)

	for i := 0; i < 100; i++ {
		key, err := adapter.GenerateKey()
		if err != nil {
			t.Fatalf("GenerateKey() iteration %d error = %v", i, err)
		}

		keyStr := string(key)
		if keys[keyStr] {
			t.Errorf("GenerateKey() produced duplicate key on iteration %d", i)
		}

		keys[keyStr] = true
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	adapter := aes.Adapter{
		Key: func() []byte { return key },
	}

	tests := []struct {
		name  string
		plain []byte
	}{
		{"empty", []byte{}},
		{"short", []byte("hello")},
		{"medium", []byte("the quick brown fox jumps over the lazy dog")},
		{"with-nulls", []byte("hello\x00world")},
		{"binary", []byte{0x00, 0x01, 0x02, 0xff, 0xfe, 0xfd}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := adapter.Encrypt(tt.plain)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			decrypted, err := adapter.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if string(decrypted) != string(tt.plain) {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.plain)
			}
		})
	}
}

func TestEncrypt_ProducesDifferentCiphertexts(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	adapter := aes.Adapter{
		Key: func() []byte { return key },
	}

	plain := []byte("same plaintext")
	ciphertexts := make(map[string]bool)

	for i := 0; i < 100; i++ {
		encrypted, err := adapter.Encrypt(plain)
		if err != nil {
			t.Fatalf("Encrypt() iteration %d error = %v", i, err)
		}

		if ciphertexts[encrypted] {
			t.Errorf("Encrypt() produced duplicate ciphertext on iteration %d", i)
		}

		ciphertexts[encrypted] = true
	}
}

func TestEncryptDecrypt_WithAAD(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	adapter := aes.Adapter{
		Key: func() []byte { return key },
	}

	plain := []byte("secret message")
	aad := []byte("additional authenticated data")

	encrypted, err := adapter.Encrypt(plain, aad)
	if err != nil {
		t.Fatalf("Encrypt() with AAD error = %v", err)
	}

	// Decrypt with correct AAD
	decrypted, err := adapter.Decrypt(encrypted, aad)
	if err != nil {
		t.Fatalf("Decrypt() with correct AAD error = %v", err)
	}

	if string(decrypted) != string(plain) {
		t.Errorf("Decrypt() = %v, want %v", decrypted, plain)
	}
}

func TestDecrypt_WrongAAD(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	adapter := aes.Adapter{
		Key: func() []byte { return key },
	}

	plain := []byte("secret message")
	aad := []byte("correct aad")

	encrypted, err := adapter.Encrypt(plain, aad)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Try to decrypt with wrong AAD
	_, err = adapter.Decrypt(encrypted, []byte("wrong aad"))
	if err == nil {
		t.Error("Decrypt() with wrong AAD should fail")
	}
}

func TestDecrypt_MissingAAD(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	adapter := aes.Adapter{
		Key: func() []byte { return key },
	}

	plain := []byte("secret message")
	aad := []byte("additional data")

	encrypted, err := adapter.Encrypt(plain, aad)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Try to decrypt without AAD when it was encrypted with AAD
	_, err = adapter.Decrypt(encrypted)
	if err == nil {
		t.Error("Decrypt() without AAD when encrypted with AAD should fail")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	adapter := aes.Adapter{
		Key: func() []byte { return key },
	}

	_, err := adapter.Decrypt("not-valid-base64!!!")
	if err == nil {
		t.Error("Decrypt() with invalid base64 should fail")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	adapter := aes.Adapter{
		Key: func() []byte { return key },
	}

	// Base64 of just a few bytes - too short to contain nonce
	_, err := adapter.Decrypt("AAAA")
	if err == nil {
		t.Error("Decrypt() with too short ciphertext should fail")
	}
}

func TestDecrypt_WithPreviousKeys(t *testing.T) {
	oldKey := make([]byte, 32)

	newKey := make([]byte, 32)
	for i := range oldKey {
		oldKey[i] = byte(i)
		newKey[i] = byte(i + 100)
	}

	// Encrypt with old key
	oldAdapter := aes.Adapter{
		Key: func() []byte { return oldKey },
	}

	plain := []byte("secret message")

	encrypted, err := oldAdapter.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Decrypt with new adapter that has old key in PreviousKeys
	newAdapter := aes.Adapter{
		Key:          func() []byte { return newKey },
		PreviousKeys: func() [][]byte { return [][]byte{oldKey} },
	}

	decrypted, err := newAdapter.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt() with previous keys error = %v", err)
	}

	if string(decrypted) != string(plain) {
		t.Errorf("Decrypt() = %v, want %v", decrypted, plain)
	}
}

func TestEncrypt_InvalidKeySize(t *testing.T) {
	adapter := aes.Adapter{
		Key: func() []byte { return []byte("short") }, // Invalid key size
	}

	_, err := adapter.Encrypt([]byte("test"))
	if err == nil {
		t.Error("Encrypt() with invalid key size should fail")
	}
}

func TestDecrypt_InvalidKeySize(t *testing.T) {
	// First encrypt with valid key
	validKey := make([]byte, 32)
	validAdapter := aes.Adapter{
		Key: func() []byte { return validKey },
	}

	encrypted, _ := validAdapter.Encrypt([]byte("test"))

	// Try to decrypt with invalid key
	invalidAdapter := aes.Adapter{
		Key: func() []byte { return []byte("short") },
	}

	_, err := invalidAdapter.Decrypt(encrypted)
	if err == nil {
		t.Error("Decrypt() with invalid key size should fail")
	}
}
