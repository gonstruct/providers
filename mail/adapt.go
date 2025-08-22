package mail

import (
	"embed"

	"github.com/gonstruct/providers/contracts"
	"github.com/gonstruct/providers/entities/mailables"
)

var globalProvider *provider

type provider struct {
	adapter         contracts.Mail
	templates       embed.FS
	defaultEnvelope *mailables.Envelope
}

func Adapt(adapter contracts.Mail, options ...func(*provider)) {
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

func WithTemplates(templates embed.FS) func(*provider) {
	return func(p *provider) {
		p.templates = templates
	}
}

func WithDefaultEnvelope(envelope mailables.Envelope) func(*provider) {
	return func(p *provider) {
		p.defaultEnvelope = &envelope
	}
}
