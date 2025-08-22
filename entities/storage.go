package entities

import (
	"strings"

	"github.com/gonstruct/providers/entities/file"
)

type StorageInput struct {
	ID   string
	File file.File
	Path string
}

func (i StorageInput) Name() string {
	return strings.TrimSuffix(i.File.Name, i.File.Extension())
}

type StorageObject struct {
	Name     string
	Path     string
	MimeType string
}
