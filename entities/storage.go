package entities

import (
	"strings"
	"time"

	"github.com/gonstruct/providers/entities/file"
)

// Visibility represents file access permissions.
type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

type StorageInput struct {
	ID   string
	File file.File
	Path string
}

func (i StorageInput) Name() string {
	return strings.TrimSuffix(i.File.Name, i.File.Extension())
}

// StorageObject contains metadata about a stored file.
type StorageObject struct {
	Name         string
	Path         string
	MimeType     string
	Size         int64
	LastModified time.Time
	Visibility   Visibility
}
