package fake

func (a *Adapter) GenerateKey() ([]byte, error) {
	if a.GenerateKeyFunc != nil {
		return a.GenerateKeyFunc()
	}

	if a.GenerateKeyError != nil {
		return nil, a.GenerateKeyError
	}

	// Return a fake 32-byte key
	return []byte("fake-encryption-key-32-bytes!!!"), nil
}
