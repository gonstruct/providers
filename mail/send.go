package mail

import (
	"github.com/gonstruct/providers/contracts"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/entities/mailables"
)

func Send(mailable contracts.Mailable, optionSlice ...Option) error {
	options := apply(optionSlice...)

	envelope := buildEnvelope(options.DefaultEnvelope, mailable.Envelope())

	view := mailable.Content()

	content, err := view.Parse(options.Templates)
	if err != nil {
		return err
	}

	text, err := view.ParseText(options.Templates)
	if err != nil {
		return err
	}

	return options.Adapter.Send(options.Context, entities.MailInput{
		Envelope:    envelope,
		Attachments: mailable.Attachments(),
		Html:        content,
		Text:        text,
	})
}

func buildEnvelope(defaultEnvelope *mailables.Envelope, override mailables.Envelope) mailables.Envelope {
	if defaultEnvelope == nil {
		return override
	}

	base := *defaultEnvelope
	base.To = append(base.To[:0:0], defaultEnvelope.To...)
	base.Cc = append(base.Cc[:0:0], defaultEnvelope.Cc...)
	base.Bcc = append(base.Bcc[:0:0], defaultEnvelope.Bcc...)
	base.Headers = append(base.Headers[:0:0], defaultEnvelope.Headers...)

	return base.Merge(override)
}
