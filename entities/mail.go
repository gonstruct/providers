package entities

import (
	"bytes"

	"github.com/gonstruct/providers/entities/mailables"
)

type MailInput struct {
	Envelope    mailables.Envelope
	Attachments mailables.AttachmentSlice
	Html        bytes.Buffer
}
