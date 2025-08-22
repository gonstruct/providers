package smtp

import (
	"bytes"
	"context"
	"errors"

	"github.com/gonstruct/providers/entities/mailables"
	"gopkg.in/gomail.v2"
)

func (adapter *Adapter) Send(context context.Context, envelope mailables.Envelope, body bytes.Buffer) error {
	message := gomail.NewMessage()

	if envelope.Subject != "" {
		message.SetHeader("Subject", envelope.Subject)
	} else {
		return errors.New("no subject specified")
	}

	if envelope.From != nil {
		message.SetHeader("From", envelope.From.String())
	} else {
		return errors.New("no sender specified")
	}

	if len(envelope.To) > 0 {
		message.SetHeader("To", envelope.To.String()...)
	} else {
		return errors.New("no recipient(s) specified")
	}

	if envelope.ReplyTo != nil {
		message.SetHeader("Reply-To", envelope.ReplyTo.String())
	}

	if len(envelope.Cc) > 0 {
		message.SetHeader("Cc", envelope.Cc.String()...)
	}

	if len(envelope.Bcc) > 0 {
		message.SetHeader("Bcc", envelope.Bcc.String()...)
	}

	message.SetBody("text/html", body.String())

	return gomail.
		NewDialer(adapter.Host, adapter.Port, adapter.Username, adapter.Password).
		DialAndSend(message)
}
