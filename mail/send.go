package mail

import "github.com/gonstruct/providers/contracts"

func Send(mailable contracts.Mailable, optionSlice ...Option) error {
	options := apply(optionSlice...)

	envelope := options.DefaultEnvelope.Merge(mailable.Envelope())

	content, err := mailable.Content().Parse(options.Templates)
	if err != nil {
		return err
	}

	return options.Adapter.Send(options.Context, envelope, content)
}
