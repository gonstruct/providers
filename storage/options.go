package storage

import (
	"context"

	"github.com/gonstruct/providers/contracts"
	"github.com/google/uuid"
)

type options struct {
	Context          context.Context
	Adapter          contracts.Storage
	GenerateUniqueID func() string
}

type Option func(*options)

func apply(optionSlice ...Option) *options {
	options := &options{
		Context:          context.Background(),
		Adapter:          globalProvider.adapter,
		GenerateUniqueID: func() string { return uuid.NewString() },
	}

	for _, option := range optionSlice {
		option(options)
	}

	return options
}

func WithContext(ctx context.Context) Option {
	return func(options *options) {
		options.Context = ctx
	}
}

func WithAdapter(adapter contracts.Storage) Option {
	return func(options *options) {
		options.Adapter = adapter
	}
}

func WithUniqueIDGenerator(generateUniqueID func() string) Option {
	return func(options *options) {
		options.GenerateUniqueID = generateUniqueID
	}
}
