package storage

import (
	"github.com/gonstruct/providers/contracts"
)

var globalProvider *provider

type provider struct {
	adapter contracts.Storage
}

func Adapt(adapter contracts.Storage, options ...func(*provider)) {
	provider := &provider{
		adapter: adapter,
	}

	for _, option := range options {
		option(provider)
	}

	if globalProvider != nil {
		panic("mail provider already set")
	}

	globalProvider = provider
}
