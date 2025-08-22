package storage

import (
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/entities/file"
)

func PutFile(path string, file file.File, optionSlice ...Option) (*entities.StorageObject, error) {
	options := apply(optionSlice...)

	return options.Adapter.PutFile(options.Context, entities.StorageInput{
		ID:   options.GenerateUniqueID(),
		File: file,
		Path: path,
	})
}
