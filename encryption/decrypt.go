package encryption

import (
	"errors"
)

// Decrypt decrypts data using the configured encryption adapter
// Optional additionalData must match the AAD used during encryption
func Decrypt[T stringOrBytes](base64Cipher T, optionSlice ...Option) ([]byte, error) {
	options := apply(optionSlice...)

	switch data := any(base64Cipher).(type) {
	case string:
		return options.Adapter.Decrypt(data, options.AdditionalData...)
	case []byte:
		return options.Adapter.Decrypt(string(data), options.AdditionalData...)
	}

	return nil, errors.New("unsupported data type")
}

// DecryptString decrypts data and returns it as a string
func DecryptString[T stringOrBytes](base64Cipher T, optionSlice ...Option) (string, error) {
	result, err := Decrypt(base64Cipher, optionSlice...)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
