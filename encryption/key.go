package encryption

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type keyFormat string

const (
	KeyFormatBase64 keyFormat = "base64"
	KeyFormatHex    keyFormat = "hex"
)

func GenerateKey(format keyFormat, optionSlice ...Option) (string, error) {
	options := apply(optionSlice...)

	key, err := options.Adapter.GenerateKey()
	if err != nil {
		return "", err
	}

	switch format {
	case KeyFormatBase64:
		return strings.Join([]string{"base64", base64.StdEncoding.EncodeToString(key)}, ":"), nil
	case KeyFormatHex:
		return strings.Join([]string{"hex", hex.EncodeToString(key)}, ":"), nil
	default:
		return "", errors.New("unsupported key format")
	}
}

func ParseKeys(keys ...string) ([][]byte, error) {
	byteKeys := make([][]byte, 0, len(keys))

	for _, key := range keys {
		parsedKey, err := ParseKey(key)
		if err != nil {
			return nil, err
		}

		byteKeys = append(byteKeys, parsedKey)
	}

	return byteKeys, nil
}

func MustParseKeys(keys ...string) [][]byte {
	byteKeys, err := ParseKeys(keys...)
	if err != nil {
		panic(err)
	}

	return byteKeys
}

func ParseKey(encoded string) ([]byte, error) {
	splitN := 2

	splitted := strings.SplitN(encoded, ":", splitN)
	if len(splitted) != splitN {
		return nil, fmt.Errorf("invalid key format: %s", encoded)
	}

	switch keyFormat(splitted[0]) {
	case KeyFormatBase64:
		return base64.StdEncoding.DecodeString(splitted[1])
	case KeyFormatHex:
		return hex.DecodeString(splitted[1])
	default:
		return nil, fmt.Errorf("unsupported key format: %s", splitted[0])
	}
}

func MustParseKey(encoded string) []byte {
	key, err := ParseKey(encoded)
	if err != nil {
		panic(err)
	}

	return key
}
