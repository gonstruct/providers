package contracts

import (
	"context"

	"github.com/gonstruct/providers/entities"
)

type Mail interface {
	Send(context context.Context, input entities.MailInput) error
}
