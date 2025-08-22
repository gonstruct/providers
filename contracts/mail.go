package contracts

import (
	"bytes"
	"context"

	"github.com/gonstruct/providers/entities/mailables"
)

type MailAdapter interface {
	Send(context context.Context, envelope mailables.Envelope, html bytes.Buffer) error
}
