package contracts

import (
	"context"

	"github.com/gonstruct/providers/entities"
)

type MailAdapter interface {
	Send(context context.Context, input entities.MailInput) error
}
