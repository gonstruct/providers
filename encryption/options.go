package encryption

import (
	"github.com/gonstruct/providers/contracts"
)

type options struct {
	Adapter contracts.Encryption
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
