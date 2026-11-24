package smtp

import (
	"context"
	"io"

	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/mail"
	"gopkg.in/gomail.v2"
)

func (adapter *Adapter) Send(ctx context.Context, input entities.MailInput) error {
	_ = ctx // gomail doesn't support context, but we accept it for interface compatibility
	message := gomail.NewMessage()

	if err := mail.CheckHeaders(input.Envelope.Headers); err != nil {
		return mail.Err("validate", err)
	}

	if input.Envelope.Subject != "" {
		message.SetHeader("Subject", input.Envelope.Subject)
	} else {
		return mail.Err("validate", mail.ErrNoSubject)
	}

	if input.Envelope.From != nil {
		message.SetHeader("From", input.Envelope.From.String())
	} else {
		return mail.Err("validate", mail.ErrNoSender)
	}

	if len(input.Envelope.To) > 0 {
		message.SetHeader("To", input.Envelope.To.String()...)
	} else {
		return mail.Err("validate", mail.ErrNoRecipients)
	}

	if input.Envelope.ReplyTo != nil {
		message.SetHeader("Reply-To", input.Envelope.ReplyTo.String())
	}

	if len(input.Envelope.Cc) > 0 {
		message.SetHeader("Cc", input.Envelope.Cc.String()...)
	}

	if len(input.Envelope.Bcc) > 0 {
		message.SetHeader("Bcc", input.Envelope.Bcc.String()...)
	}

	setExtraHeaders(message, input)

	if len(input.Attachments) > 0 {
		for _, attachment := range input.Attachments {
			message.Attach(attachment.Name,
				gomail.SetHeader(map[string][]string{
					"Content-Type": {attachment.Mime},
				}),
				gomail.SetCopyFunc(func(w io.Writer) error {
					_, err := w.Write(attachment.Content())

					return err
				}),
			)
		}
	}

	setBody(message, input)

	dialer := gomail.NewDialer(adapter.Host, adapter.Port, adapter.Username, adapter.Password)
	if err := dialer.DialAndSend(message); err != nil {
		return mail.Err("send via SMTP", err)
	}

	return nil
}

// setExtraHeaders sets the envelope's own headers, and its SES configuration
// set as the header SES's SMTP interface reads.
func setExtraHeaders(message *gomail.Message, input entities.MailInput) {
	for _, header := range input.Envelope.Headers {
		message.SetHeader(header.Name, header.Value)
	}

	if input.Envelope.ConfigurationSet != "" {
		message.SetHeader("X-SES-CONFIGURATION-SET", input.Envelope.ConfigurationSet)
	}
}

// setBody sets the HTML, and with a text version makes the HTML its
// alternative, last, as the one a client prefers when it can show it.
func setBody(message *gomail.Message, input entities.MailInput) {
	if input.Text.Len() > 0 {
		message.SetBody("text/plain", input.Text.String())
		message.AddAlternative("text/html", input.Html.String())

		return
	}

	message.SetBody("text/html", input.Html.String())
}
