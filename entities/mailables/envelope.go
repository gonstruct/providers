package mailables

type Envelope struct {
	From    *address
	To      addressSlice
	Cc      addressSlice
	Bcc     addressSlice
	ReplyTo *address
	Subject string
	// Headers are added to the message as they are, in order: a
	// List-Unsubscribe, an X-Entity-Ref-ID. The ones the envelope already
	// sets (From, To, Subject and the like) are not for here.
	Headers []Header
	// ConfigurationSet is the Amazon SES configuration set the mail goes out
	// with, which decides where its bounces and complaints are reported.
	// Over SMTP it travels as X-SES-CONFIGURATION-SET, which SES's SMTP
	// interface reads and any other server ignores.
	ConfigurationSet string
	// Tags     []string
	// Metadata map[string]any
}

// Header is one header of a message.
type Header struct {
	Name  string
	Value string
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

	if len(override.Headers) > 0 {
		envelope.Headers = append(envelope.Headers, override.Headers...)
	}

	if override.ConfigurationSet != "" {
		envelope.ConfigurationSet = override.ConfigurationSet
	}

	return *envelope
}
