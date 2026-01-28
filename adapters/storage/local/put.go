package local

import (
	"context"
	"io"
	"os"
	"path/filepath"

	gomime "github.com/cubewise-code/go-mime"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/storage"
)

// PutFile stores a file with a unique ID
func (a *Adapter) PutFile(ctx context.Context, input entities.StorageInput) (*entities.StorageObject, error) {
	extension := input.File.Extension()
	mimetype := gomime.TypeByExtension(extension)
	key := filepath.Join(input.Path, input.ID+extension)

	fullPath := filepath.Join(a.Root, key)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), os.FileMode(a.DirectoryPermission)); err != nil {
		return nil, storage.PathErr("create directory", filepath.Dir(key), err)
	}

	// Read all content from the file body
	content, err := io.ReadAll(input.File.Body)
	if err != nil {
		return nil, storage.Err("read file body", err)
	}

	// Write file
	if err := os.WriteFile(fullPath, content, os.FileMode(a.FilePermission)); err != nil {
		return nil, storage.PathErr("write file", key, err)
	}

	return &entities.StorageObject{
		Name:     input.Name(),
		Path:     key,
		MimeType: mimetype,
	}, nil
}

// Put stores raw bytes at the given path
func (a *Adapter) Put(ctx context.Context, path string, contents []byte) error {
	fullPath := filepath.Join(a.Root, path)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), os.FileMode(a.DirectoryPermission)); err != nil {
		return storage.PathErr("create directory", filepath.Dir(path), err)
	}

	if err := os.WriteFile(fullPath, contents, os.FileMode(a.FilePermission)); err != nil {
		return storage.PathErr("write file", path, err)
	}

	return nil
}

// PutStream stores content from a reader at the given path
func (a *Adapter) PutStream(ctx context.Context, path string, stream io.Reader) error {
	fullPath := filepath.Join(a.Root, path)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), os.FileMode(a.DirectoryPermission)); err != nil {
		return storage.PathErr("create directory", filepath.Dir(path), err)
	}

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(a.FilePermission))
	if err != nil {
		return storage.PathErr("create file", path, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, stream); err != nil {
		return storage.PathErr("write stream", path, err)
	}

	return nil
}
