package contracts

import (
	"context"

	"github.com/gonstruct/providers/entities"
)

type Storage interface {
	PutFile(context context.Context, input entities.StorageInput) (*entities.StorageObject, error)
}
