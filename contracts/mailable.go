package contracts

import "github.com/gonstruct/providers/entities/mailables"

type Mailable interface {
	Envelope() mailables.Envelope
	Content() mailables.Content
	// Attachments() mailables.AttachmentSlice
}
