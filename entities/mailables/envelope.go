package mailables

type Envelope struct {
	From    *address
	To      addressSlice
	Cc      addressSlice
	Bcc     addressSlice
	ReplyTo *address
	Subject string
	// Tags     []string
	// Metadata map[string]any
}

func (envelope *Envelope) Merge(override Envelope) Envelope {
	if override.From != nil {
		envelope.From = override.From
	}

	if len(override.To) > 0 {
		envelope.To = append(envelope.To, override.To...)
	}

	if len(override.Cc) > 0 {
		envelope.Cc = append(envelope.Cc, override.Cc...)
	}

	if len(override.Bcc) > 0 {
		envelope.Bcc = append(envelope.Bcc, override.Bcc...)
	}

	if override.ReplyTo != nil {
		envelope.ReplyTo = override.ReplyTo
	}

	if override.Subject != "" {
		envelope.Subject = override.Subject
	}

	return *envelope
}
