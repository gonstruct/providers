package local

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/storage"
)

// GetVisibility returns the visibility of a file based on its permissions.
func (a *Adapter) GetVisibility(ctx context.Context, path string) (entities.Visibility, error) {
	fullPath := filepath.Join(a.Root, path)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", storage.PathErr("get visibility", path, storage.ErrFileNotFound)
		}

		return "", storage.PathErr("get visibility", path, err)
	}

	mode := info.Mode().Perm()
	// Check if others have read permission (public)
	if mode&0o004 != 0 {
		return entities.VisibilityPublic, nil
	}

	return entities.VisibilityPrivate, nil
}

// SetVisibility changes the visibility of a file.
func (a *Adapter) SetVisibility(ctx context.Context, path string, visibility entities.Visibility) error {
	fullPath := filepath.Join(a.Root, path)

	perm := a.getVisibilityPermission(visibility)
	if err := os.Chmod(fullPath, os.FileMode(perm)); err != nil {
		if os.IsNotExist(err) {
			return storage.PathErr("set visibility", path, storage.ErrFileNotFound)
		}

		return storage.PathErr("set visibility", path, err)
	}

	return nil
}
