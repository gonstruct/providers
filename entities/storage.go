package entities

import (
	"net/http"
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
	Name     string
	Path     string
	MimeType string
}

// TemporaryStorageObject contains information about a temporary storage object, such as a presigned URL.
type TemporaryStorageInput struct {
	ID         string
	File       file.File
	Path       string
	Visibility Visibility
	Expiry     time.Duration
}

func (i TemporaryStorageInput) Name() string {
	return strings.TrimSuffix(i.File.Name, i.File.Extension())
}

type TemporaryStorageObject struct {
	Name         string
	Path         string
	URL          string
	Method       string
	SignedHeader http.Header
}
