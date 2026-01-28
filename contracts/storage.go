package contracts

import (
	"context"
	"io"
	"time"

	"github.com/gonstruct/providers/entities"
)

// Storage defines the interface for storage operations (Laravel-style).
type Storage interface {
	// Writing files
	PutFile(ctx context.Context, input entities.StorageInput) (*entities.StorageObject, error)
	Put(ctx context.Context, path string, contents []byte) error
	PutStream(ctx context.Context, path string, stream io.Reader) error

	// Reading files
	Get(ctx context.Context, path string) ([]byte, error)
	GetStream(ctx context.Context, path string) (io.ReadCloser, error)
	Exists(ctx context.Context, path string) (bool, error)
	Missing(ctx context.Context, path string) (bool, error)

	// File metadata
	Size(ctx context.Context, path string) (int64, error)
	LastModified(ctx context.Context, path string) (time.Time, error)
	MimeType(ctx context.Context, path string) (string, error)

	// File operations
	Copy(ctx context.Context, from, to string) (*entities.StorageObject, error)
	Move(ctx context.Context, from, to string) (*entities.StorageObject, error)
	Delete(ctx context.Context, paths ...string) error

	// Visibility
	GetVisibility(ctx context.Context, path string) (entities.Visibility, error)
	SetVisibility(ctx context.Context, path string, visibility entities.Visibility) error

	// Directories
	Files(ctx context.Context, directory string) ([]string, error)
	AllFiles(ctx context.Context, directory string) ([]string, error)
	Directories(ctx context.Context, directory string) ([]string, error)
	AllDirectories(ctx context.Context, directory string) ([]string, error)
	MakeDirectory(ctx context.Context, path string) error
	DeleteDirectory(ctx context.Context, directory string) error

	// URLs (optional - may return empty string if not supported)
	URL(path string) string
	TemporaryURL(ctx context.Context, path string, expiration time.Duration) (string, error)
}
