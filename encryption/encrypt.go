package encryption

import (
	"errors"
)

type stringOrBytes interface {
	~string | ~[]byte
}

// Encrypt encrypts data using the configured encryption adapter
// Optional additionalData provides authenticated data (AAD) for AES-GCM
func Encrypt[T stringOrBytes](plain T, optionSlice ...Option) (string, error) {
	options := apply(optionSlice...)

	switch data := any(plain).(type) {
	case string:
		return options.Adapter.Encrypt([]byte(data), options.AdditionalData...)
	case []byte:
		return options.Adapter.Encrypt(data, options.AdditionalData...)
	}

	return "", errors.New("unsupported data type")
}
