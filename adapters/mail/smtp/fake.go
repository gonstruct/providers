package smtp

import (
	"context"

	"github.com/gonstruct/providers/entities"
)

// FakeAdapter is a mock SMTP adapter for testing
type FakeAdapter struct {
	// SendFunc allows customizing the Send behavior
	SendFunc func(ctx context.Context, input entities.MailInput) error

	// SendCalls records all calls to Send
	SendCalls []FakeSendCall
}

type FakeSendCall struct {
	Context context.Context
	Input   entities.MailInput
}

// Fake creates a new mock SMTP adapter with default behaviors
func Fake() *FakeAdapter {
	return &FakeAdapter{
		SendFunc: func(ctx context.Context, input entities.MailInput) error {
			return nil
		},
	}
}

func (a *FakeAdapter) Send(ctx context.Context, input entities.MailInput) error {
	a.SendCalls = append(a.SendCalls, FakeSendCall{
		Context: ctx,
		Input:   input,
	})
	return a.SendFunc(ctx, input)
}

// Reset clears all recorded calls
func (a *FakeAdapter) Reset() {
	a.SendCalls = nil
}

// LastSendCall returns the most recent Send call, or nil if none
func (a *FakeAdapter) LastSendCall() *FakeSendCall {
	if len(a.SendCalls) == 0 {
		return nil
	}
	return &a.SendCalls[len(a.SendCalls)-1]
}
