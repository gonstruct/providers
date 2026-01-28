package local

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/gonstruct/providers/storage"
)

// Copy copies a file from one location to another
func (a *Adapter) Copy(ctx context.Context, from, to string) error {
	srcPath := filepath.Join(a.Root, from)
	dstPath := filepath.Join(a.Root, to)

	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return storage.PathErr("copy", from, storage.ErrFileNotFound)
		}
		return storage.PathErr("copy open source", from, err)
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return storage.PathErr("copy stat source", from, err)
	}

	// Create destination directory if needed
	if err := os.MkdirAll(filepath.Dir(dstPath), os.FileMode(a.DirectoryPermission)); err != nil {
		return storage.PathErr("copy create directory", filepath.Dir(to), err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return storage.PathErr("copy create destination", to, err)
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return storage.Err("copy contents", err)
	}

	return nil
}

// Move moves a file from one location to another
func (a *Adapter) Move(ctx context.Context, from, to string) error {
	srcPath := filepath.Join(a.Root, from)
	dstPath := filepath.Join(a.Root, to)

	// Create destination directory if needed
	if err := os.MkdirAll(filepath.Dir(dstPath), os.FileMode(a.DirectoryPermission)); err != nil {
		return storage.PathErr("move create directory", filepath.Dir(to), err)
	}

	// Try rename first (works if same filesystem)
	if err := os.Rename(srcPath, dstPath); err != nil {
		// Fall back to copy+delete for cross-filesystem moves
		if err := a.Copy(ctx, from, to); err != nil {
			return err
		}
		if err := os.Remove(srcPath); err != nil {
			return storage.PathErr("move remove source", from, err)
		}
	}

	return nil
}

// Delete removes one or more files
func (a *Adapter) Delete(ctx context.Context, paths ...string) error {
	for _, path := range paths {
		fullPath := filepath.Join(a.Root, path)
		if err := os.Remove(fullPath); err != nil {
			if os.IsNotExist(err) {
				continue // Ignore non-existent files (idempotent delete)
			}
			return storage.PathErr("delete", path, err)
		}
	}

	return nil
}
