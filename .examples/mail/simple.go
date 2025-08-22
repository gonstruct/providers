package mail

import (
	"context"
	"embed"

	"github.com/gonstruct/providers/adapters/mail/amazon_ses"
	"github.com/gonstruct/providers/adapters/mail/smtp"
	"github.com/gonstruct/providers/facades"
	"github.com/gonstruct/providers/providers/mail"
	"github.com/gonstruct/providers/providers/mail/mailables"
)

func VerificationMail(data any) *facades.Mailable[verificationMail] {
	return facades.MailableOf(verificationMail{
		Data: data,
	})
}

type verificationMail struct {
	Data any
}

func (verificationMail) Envelope() mailables.Envelope {
	return mailables.Envelope{
		From: mailables.Address("email@email.ts", "hello world"),
	}
}

func (v verificationMail) Content() mailables.Content {
	return mailables.Content{
		View: "emails/verification.html",
		With: map[string]any{
			"data": v.Data,
		},
	}
}

func ExampleSimple() {
	ctx := context.Background()

	mailable := VerificationMail("hello world")
	_ = mailable.Send()

	_ = mailable.Send(
		mail.WithContext(ctx),
		mail.WithAdapter(&amazon_ses.Adapter{}),
	)
}

var templates embed.FS

func init() {
	mail.Provide(
		&smtp.Adapter{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "username",
			Password: "password",
		},
		mail.WithTemplates(templates),
		mail.WithDefaultEnvelope(mailables.Envelope{
			From: mailables.Address("default@company.com", "Default Company"),
		}),
	)
}
