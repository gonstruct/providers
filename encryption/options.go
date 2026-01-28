package encryption

import (
	"github.com/gonstruct/providers/contracts"
)

type options struct {
	Adapter        contracts.Encryption
	AdditionalData [][]byte
}

type Option func(*options)

func apply(optionSlice ...Option) *options {
	options := &options{}

	if globalProvider != nil {
		options.Adapter = globalProvider.adapter
	}

	for _, option := range optionSlice {
		option(options)
	}

	return options
}

func WithAdapter(adapter contracts.Encryption) Option {
	return func(options *options) {
		options.Adapter = adapter
	}
}

// WithAdditionalData sets the Additional Authenticated Data (AAD) for AES-GCM
// AAD is authenticated but not encrypted - useful for binding ciphertext to context
// Example use cases: user ID, record ID, timestamp, etc.
func WithAdditionalData(aad []byte) Option {
	return func(options *options) {
		options.AdditionalData = [][]byte{aad}
	}
}
