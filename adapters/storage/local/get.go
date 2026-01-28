package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	gomime "github.com/cubewise-code/go-mime"
	"github.com/gonstruct/providers/storage"
)

// Get retrieves the contents of a file
func (a *Adapter) Get(ctx context.Context, path string) ([]byte, error) {
	fullPath := filepath.Join(a.Root, path)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.PathErr("get", path, storage.ErrFileNotFound)
		}
		return nil, storage.PathErr("get", path, err)
	}

	return content, nil
}

// GetStream returns a reader for the file contents
func (a *Adapter) GetStream(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(a.Root, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.PathErr("open stream", path, storage.ErrFileNotFound)
		}
		return nil, storage.PathErr("open stream", path, err)
	}

	return file, nil
}

// Exists checks if a file exists
func (a *Adapter) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(a.Root, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, storage.PathErr("stat", path, err)
	}

	return true, nil
}

// Missing checks if a file does not exist
func (a *Adapter) Missing(ctx context.Context, path string) (bool, error) {
	exists, err := a.Exists(ctx, path)
	return !exists, err
}

// Size returns the size of a file in bytes
func (a *Adapter) Size(ctx context.Context, path string) (int64, error) {
	fullPath := filepath.Join(a.Root, path)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, storage.PathErr("size", path, storage.ErrFileNotFound)
		}
		return 0, storage.PathErr("size", path, err)
	}

	return info.Size(), nil
}

// LastModified returns the last modification time of a file
func (a *Adapter) LastModified(ctx context.Context, path string) (time.Time, error) {
	fullPath := filepath.Join(a.Root, path)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Time{}, storage.PathErr("last modified", path, storage.ErrFileNotFound)
		}
		return time.Time{}, storage.PathErr("last modified", path, err)
	}

	return info.ModTime(), nil
}

// MimeType returns the MIME type of a file based on its extension
func (a *Adapter) MimeType(ctx context.Context, path string) (string, error) {
	// First check if file exists
	exists, err := a.Exists(ctx, path)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", storage.PathErr("mime type", path, storage.ErrFileNotFound)
	}

	ext := filepath.Ext(path)
	mimeType := gomime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return mimeType, nil
}
