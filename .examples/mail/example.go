package mail

import (
	"net/mail"

	"github.com/gonstruct/providers/adapters/mail/smtp"
	"github.com/gonstruct/providers/entities/mailables"
	pmail "github.com/gonstruct/providers/mail"
)

// WelcomeMail implements contracts.Mailable
type WelcomeMail struct {
	To   string
	Name string
}

func (m WelcomeMail) Envelope() mailables.Envelope {
	return mailables.Envelope{
		Subject: "Welcome!",
		To:      []*mail.Address{{Address: m.To, Name: m.Name}},
	}
}

func (m WelcomeMail) Content() mailables.Content {
	return mailables.Content{
		View: "emails/welcome.html",
		With: map[string]any{"name": m.Name},
	}
}

func (m WelcomeMail) Attachments() mailables.AttachmentSlice { return nil }

func Example() {
	pmail.Adapt(
		&smtp.Adapter{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "user",
			Password: "pass",
		},
		pmail.WithDefaultEnvelope(mailables.Envelope{
			From: mailables.Address("noreply@app.com", "My App"),
		}),
	)

	_ = pmail.Send(WelcomeMail{To: "john@example.com", Name: "John"})
}
