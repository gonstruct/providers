package mail

import (
	"context"
	"embed"

	"github.com/gonstruct/providers/contracts"
	"github.com/gonstruct/providers/entities/mailables"
)

type options struct {
	Context         context.Context
	Adapter         contracts.Mail
	Templates       embed.FS
	DefaultEnvelope *mailables.Envelope
}

type Option func(*options)

func apply(optionSlice ...Option) *options {
	options := &options{
		Context:         context.Background(),
		Adapter:         globalProvider.adapter,
		Templates:       globalProvider.templates,
		DefaultEnvelope: globalProvider.defaultEnvelope,
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

func WithAdapter(adapter contracts.Mail) Option {
	return func(options *options) {
		options.Adapter = adapter
	}
}
