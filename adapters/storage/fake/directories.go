package fake

import (
	"context"
)

func (a *Adapter) Files(ctx context.Context, directory string) ([]string, error) {
	if a.FilesError != nil {
		return nil, a.FilesError
	}

	return a.listFiles(directory, false), nil
}

func (a *Adapter) AllFiles(ctx context.Context, directory string) ([]string, error) {
	if a.FilesError != nil {
		return nil, a.FilesError
	}

	return a.listFiles(directory, true), nil
}

func (a *Adapter) listFiles(directory string, recursive bool) []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var files []string
	prefix := directory
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	for path := range a.files {
		if prefix == "" || (len(path) > len(prefix) && path[:len(prefix)] == prefix) {
			if prefix == "" {
				files = append(files, path)
			} else {
				rest := path[len(prefix):]
				if recursive || !containsSlash(rest) {
					files = append(files, path)
				}
			}
		}
	}

	return files
}

func containsSlash(s string) bool {
	for _, c := range s {
		if c == '/' {
			return true
		}
	}
	return false
}

func (a *Adapter) Directories(ctx context.Context, directory string) ([]string, error) {
	if a.DirectoriesError != nil {
		return nil, a.DirectoriesError
	}

	return a.listDirectories(directory, false), nil
}

func (a *Adapter) AllDirectories(ctx context.Context, directory string) ([]string, error) {
	if a.DirectoriesError != nil {
		return nil, a.DirectoriesError
	}

	return a.listDirectories(directory, true), nil
}

func (a *Adapter) listDirectories(directory string, recursive bool) []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	dirs := make(map[string]bool)
	prefix := directory
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	for path := range a.files {
		if prefix == "" || (len(path) > len(prefix) && path[:len(prefix)] == prefix) {
			var rest string
			if prefix == "" {
				rest = path
			} else {
				rest = path[len(prefix):]
			}
			for i, c := range rest {
				if c == '/' {
					var dir string
					if prefix == "" {
						dir = rest[:i]
					} else {
						dir = prefix + rest[:i]
					}
					dirs[dir] = true
					if !recursive {
						break
					}
				}
			}
		}
	}

	result := make([]string, 0, len(dirs))
	for dir := range dirs {
		result = append(result, dir)
	}

	return result
}

func (a *Adapter) MakeDirectory(ctx context.Context, path string) error {
	if a.MakeDirectoryError != nil {
		return a.MakeDirectoryError
	}
	return nil
}

func (a *Adapter) DeleteDirectory(ctx context.Context, directory string) error {
	if a.DeleteDirError != nil {
		return a.DeleteDirError
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	prefix := directory
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	for path := range a.files {
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			delete(a.files, path)
		}
	}

	return nil
}
