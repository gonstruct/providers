package mail

import (
	"embed"

	"github.com/gonstruct/providers/adapters/mail/fake"
	"github.com/gonstruct/providers/entities/mailables"
)

// FakeOption configures the fake mail adapter
type FakeOption func(*provider)

// WithFakeTemplates sets the templates for the fake adapter
func WithFakeTemplates(templates embed.FS) FakeOption {
	return func(p *provider) {
		p.templates = templates
	}
}

// WithFakeDefaultEnvelope sets the default envelope for the fake adapter
func WithFakeDefaultEnvelope(envelope mailables.Envelope) FakeOption {
	return func(p *provider) {
		p.defaultEnvelope = &envelope
	}
}

// Fake sets up a fake mail adapter for testing and returns it for assertions.
// This replaces any existing mail provider.
//
// Example:
//
//	func TestSendWelcomeEmail(t *testing.T) {
//	    fake := mail.Fake()
//
//	    // Your code that uses mail.Send()
//	    mail.Send(ctx, welcomeMail)
//
//	    // Assert
//	    fake.AssertSent(t)
//	    fake.AssertSentTo(t, "user@example.com")
//	}
func Fake(options ...FakeOption) *fake.Adapter {
	adapter := fake.New()

	globalProvider = &provider{
		adapter: adapter,
	}

	for _, opt := range options {
		opt(globalProvider)
	}

	return adapter
}
