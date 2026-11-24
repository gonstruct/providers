package mailables

import "strings"

type Envelope struct {
	From    *address
	To      addressSlice
	Cc      addressSlice
	Bcc     addressSlice
	ReplyTo *address
	Subject string
	// Headers are added to the message as they are, in order: a
	// List-Unsubscribe, an X-Entity-Ref-ID. The ones the envelope already
	// sets (From, To, Subject and the like) are refused when sending, as
	// is a name or value that would break the message (mail.CheckHeaders).
	Headers []Header
	// ConfigurationSet is the Amazon SES configuration set the mail goes out
	// with, which decides where its bounces and complaints are reported.
	// Over SMTP it travels as X-SES-CONFIGURATION-SET, which SES's SMTP
	// interface reads and any other server ignores.
	ConfigurationSet string
	// Tags     []string
	// Metadata map[string]any
}

// setHeader replaces the header of the same name, whatever its case, or
// adds it: a mail's List-Unsubscribe takes the default's place rather than
// going out beside it.
func setHeader(headers []Header, header Header) []Header {
	for index := range headers {
		if strings.EqualFold(headers[index].Name, header.Name) {
			headers[index] = header

			return headers
		}
	}

	return append(headers, header)
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

	for _, header := range override.Headers {
		envelope.Headers = setHeader(envelope.Headers, header)
	}

	if override.ConfigurationSet != "" {
		envelope.ConfigurationSet = override.ConfigurationSet
	}

	return *envelope
}
