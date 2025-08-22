package facades

import (
	"github.com/gonstruct/providers/contracts"
	"github.com/gonstruct/providers/mail"
)

type Mailable[T contracts.Mailable] struct {
	mailable T
}

func MailableOf[T contracts.Mailable](mailable T) *Mailable[T] {
	return &Mailable[T]{mailable: mailable}
}

func (m Mailable[T]) Send(optionSlice ...mail.Option) error {
	return mail.Send(m.mailable, optionSlice...)
}
