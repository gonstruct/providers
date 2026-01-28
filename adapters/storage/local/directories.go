package local

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gonstruct/providers/storage"
)

// Files returns a list of files in the given directory (non-recursive)
func (a *Adapter) Files(ctx context.Context, directory string) ([]string, error) {
	fullPath := filepath.Join(a.Root, directory)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, storage.PathErr("list files", directory, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(directory, entry.Name()))
		}
	}

	return files, nil
}

// AllFiles returns a list of all files in the directory and subdirectories
func (a *Adapter) AllFiles(ctx context.Context, directory string) ([]string, error) {
	fullPath := filepath.Join(a.Root, directory)

	var files []string
	err := filepath.WalkDir(fullPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			relPath, _ := filepath.Rel(a.Root, path)
			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, storage.PathErr("list all files", directory, err)
	}

	return files, nil
}

// Directories returns a list of directories in the given directory (non-recursive)
func (a *Adapter) Directories(ctx context.Context, directory string) ([]string, error) {
	fullPath := filepath.Join(a.Root, directory)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, storage.PathErr("list directories", directory, err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(directory, entry.Name()))
		}
	}

	return dirs, nil
}

// AllDirectories returns a list of all directories in the directory and subdirectories
func (a *Adapter) AllDirectories(ctx context.Context, directory string) ([]string, error) {
	fullPath := filepath.Join(a.Root, directory)

	var dirs []string
	err := filepath.WalkDir(fullPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != fullPath {
			relPath, _ := filepath.Rel(a.Root, path)
			dirs = append(dirs, relPath)
		}
		return nil
	})

	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, storage.PathErr("list all directories", directory, err)
	}

	return dirs, nil
}

// MakeDirectory creates a directory
func (a *Adapter) MakeDirectory(ctx context.Context, path string) error {
	fullPath := filepath.Join(a.Root, path)

	if err := os.MkdirAll(fullPath, os.FileMode(a.DirectoryPermission)); err != nil {
		return storage.PathErr("create directory", path, err)
	}

	return nil
}

// DeleteDirectory removes a directory and all its contents
func (a *Adapter) DeleteDirectory(ctx context.Context, directory string) error {
	fullPath := filepath.Join(a.Root, directory)

	if err := os.RemoveAll(fullPath); err != nil {
		return storage.PathErr("delete directory", directory, err)
	}

	return nil
}
