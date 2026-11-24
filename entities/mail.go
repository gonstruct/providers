package entities

import (
	"bytes"

	"github.com/gonstruct/providers/entities/mailables"
)

type MailInput struct {
	Envelope    mailables.Envelope
	Attachments mailables.AttachmentSlice
	Html        bytes.Buffer
	// Text is the plain text version, empty when the mail has none.
	Text bytes.Buffer
}
