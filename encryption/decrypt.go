package encryption

import (
	"errors"
)

func Decrypt[T stringOrBytes](base64Cipher T, optionSlice ...Option) ([]byte, error) {
	options := apply(optionSlice...)

	switch data := any(base64Cipher).(type) {
	case string:
		return options.Adapter.Decrypt(data)
	case []byte:
		return options.Adapter.Decrypt(string(data))
	}

	return nil, errors.New("unsupported data type")
}

func DecryptString[T stringOrBytes](base64Cipher T, optionSlice ...Option) (string, error) {
	result, err := Decrypt(base64Cipher, optionSlice...)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
