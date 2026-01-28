package fake

import (
	"context"

	"github.com/gonstruct/providers/entities"
)

func (a *Adapter) Send(ctx context.Context, input entities.MailInput) error {
	if a.SendFunc != nil {
		return a.SendFunc(ctx, input)
	}

	if a.SendError != nil {
		return a.SendError
	}

	envelope := input.Envelope

	to := make([]string, len(envelope.To))
	for i, addr := range envelope.To {
		to[i] = addr.Address
	}

	var from string
	if envelope.From != nil {
		from = envelope.From.Address
	}

	call := SendCall{
		To:          to,
		From:        from,
		Subject:     envelope.Subject,
		HTML:        input.Html.String(),
		Attachments: len(input.Attachments),
		Input:       input,
	}

	a.mu.Lock()
	a.Calls = append(a.Calls, call)
	a.mu.Unlock()

	return nil
}

// Ensure Adapter implements the interface.
var _ interface {
	Send(ctx context.Context, input entities.MailInput) error
} = (*Adapter)(nil)
