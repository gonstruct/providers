package encryption

import (
	"github.com/gonstruct/providers/contracts"
)

var globalProvider *provider

type provider struct {
	adapter contracts.Encryption
}

func Adapt(adapter contracts.Encryption, options ...func(*provider)) {
	provider := &provider{
		adapter: adapter,
	}

	for _, option := range options {
		option(provider)
	}

	if globalProvider != nil {
		panic("encryption provider already set")
	}

	globalProvider = provider
}
