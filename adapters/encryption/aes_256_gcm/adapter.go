package aes_256_gcm

type Adapter struct {
	Key          func() []byte
	PreviousKeys func() [][]byte
}

func (encrypter Adapter) Keys() [][]byte {
	keys := [][]byte{encrypter.Key()}
	keys = append(keys, encrypter.PreviousKeys()...)

	return keys
}
