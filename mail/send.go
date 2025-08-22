package mail

import (
	"github.com/gonstruct/providers/contracts"
	"github.com/gonstruct/providers/entities"
)

func Send(mailable contracts.Mailable, optionSlice ...Option) error {
	options := apply(optionSlice...)

	envelope := options.DefaultEnvelope.Merge(mailable.Envelope())

	content, err := mailable.Content().Parse(options.Templates)
	if err != nil {
		return err
	}

	return options.Adapter.Send(options.Context, entities.MailInput{
		Envelope:    envelope,
		Attachments: mailable.Attachments(),
		Html:        content,
	})
}
