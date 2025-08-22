package encryption

import (
	"errors"
)

type stringOrBytes interface {
	~string | ~[]byte
}

func Encrypt[T stringOrBytes](plain T, optionSlice ...Option) (string, error) {
	options := apply(optionSlice...)

	switch data := any(plain).(type) {
	case string:
		return options.Adapter.Encrypt([]byte(data))
	case []byte:
		return options.Adapter.Encrypt(data)
	}

	return "", errors.New("unsupported data type")
}
