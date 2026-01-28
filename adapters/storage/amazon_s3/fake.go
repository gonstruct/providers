package amazon_s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/gonstruct/providers/entities"
)

// FakeAdapter is an in-memory S3 adapter for testing
type FakeAdapter struct {
	mu    sync.RWMutex
	files map[string]*fakeFile

	// Error injection
	PutFileError       error
	PutError           error
	PutStreamError     error
	GetError           error
	GetStreamError     error
	ExistsError        error
	SizeError          error
	LastModifiedError  error
	MimeTypeError      error
	CopyError          error
	MoveError          error
	DeleteError        error
	GetVisibilityError error
	SetVisibilityError error
	FilesError         error
	DirectoriesError   error
	MakeDirectoryError error
	DeleteDirError     error
	TemporaryURLError  error

	BaseURL string
}

type fakeFile struct {
	Content      []byte
	MimeType     string
	Visibility   entities.Visibility
	LastModified time.Time
}

// Fake creates a new in-memory S3 adapter for testing
func Fake() *FakeAdapter {
	return &FakeAdapter{
		files:   make(map[string]*fakeFile),
		BaseURL: "https://fake-bucket.s3.amazonaws.com",
	}
}

func (a *FakeAdapter) PutFile(ctx context.Context, input entities.StorageInput) (*entities.StorageObject, error) {
	if a.PutFileError != nil {
		return nil, a.PutFileError
	}

	content, err := io.ReadAll(input.File.Body)
	if err != nil {
		return nil, err
	}

	path := input.Path + "/" + input.ID + input.File.Extension()
	now := time.Now()

	a.mu.Lock()
	a.files[path] = &fakeFile{
		Content:      content,
		MimeType:     "application/octet-stream",
		Visibility:   entities.VisibilityPrivate,
		LastModified: now,
	}
	a.mu.Unlock()

	return &entities.StorageObject{
		Name:         input.Name(),
		Path:         path,
		MimeType:     "application/octet-stream",
		Size:         int64(len(content)),
		LastModified: now,
		Visibility:   entities.VisibilityPrivate,
	}, nil
}

func (a *FakeAdapter) Put(ctx context.Context, path string, contents []byte) error {
	if a.PutError != nil {
		return a.PutError
	}

	a.mu.Lock()
	a.files[path] = &fakeFile{
		Content:      contents,
		MimeType:     "application/octet-stream",
		Visibility:   entities.VisibilityPrivate,
		LastModified: time.Now(),
	}
	a.mu.Unlock()

	return nil
}

func (a *FakeAdapter) PutStream(ctx context.Context, path string, stream io.Reader) error {
	if a.PutStreamError != nil {
		return a.PutStreamError
	}

	content, err := io.ReadAll(stream)
	if err != nil {
		return err
	}

	a.mu.Lock()
	a.files[path] = &fakeFile{
		Content:      content,
		MimeType:     "application/octet-stream",
		Visibility:   entities.VisibilityPrivate,
		LastModified: time.Now(),
	}
	a.mu.Unlock()

	return nil
}

func (a *FakeAdapter) Get(ctx context.Context, path string) ([]byte, error) {
	if a.GetError != nil {
		return nil, a.GetError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return nil, ErrFakeFileNotFound
	}

	return f.Content, nil
}

func (a *FakeAdapter) GetStream(ctx context.Context, path string) (io.ReadCloser, error) {
	if a.GetStreamError != nil {
		return nil, a.GetStreamError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return nil, ErrFakeFileNotFound
	}

	return io.NopCloser(bytes.NewReader(f.Content)), nil
}

func (a *FakeAdapter) Exists(ctx context.Context, path string) (bool, error) {
	if a.ExistsError != nil {
		return false, a.ExistsError
	}

	a.mu.RLock()
	_, ok := a.files[path]
	a.mu.RUnlock()

	return ok, nil
}

func (a *FakeAdapter) Missing(ctx context.Context, path string) (bool, error) {
	exists, err := a.Exists(ctx, path)
	return !exists, err
}

func (a *FakeAdapter) Size(ctx context.Context, path string) (int64, error) {
	if a.SizeError != nil {
		return 0, a.SizeError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return 0, ErrFakeFileNotFound
	}

	return int64(len(f.Content)), nil
}

func (a *FakeAdapter) LastModified(ctx context.Context, path string) (time.Time, error) {
	if a.LastModifiedError != nil {
		return time.Time{}, a.LastModifiedError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return time.Time{}, ErrFakeFileNotFound
	}

	return f.LastModified, nil
}

func (a *FakeAdapter) MimeType(ctx context.Context, path string) (string, error) {
	if a.MimeTypeError != nil {
		return "", a.MimeTypeError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return "", ErrFakeFileNotFound
	}

	return f.MimeType, nil
}

func (a *FakeAdapter) Copy(ctx context.Context, from, to string) error {
	if a.CopyError != nil {
		return a.CopyError
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	f, ok := a.files[from]
	if !ok {
		return ErrFakeFileNotFound
	}

	a.files[to] = &fakeFile{
		Content:      append([]byte(nil), f.Content...),
		MimeType:     f.MimeType,
		Visibility:   f.Visibility,
		LastModified: time.Now(),
	}

	return nil
}

func (a *FakeAdapter) Move(ctx context.Context, from, to string) error {
	if a.MoveError != nil {
		return a.MoveError
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	f, ok := a.files[from]
	if !ok {
		return ErrFakeFileNotFound
	}

	a.files[to] = f
	delete(a.files, from)

	return nil
}

func (a *FakeAdapter) Delete(ctx context.Context, paths ...string) error {
	if a.DeleteError != nil {
		return a.DeleteError
	}

	a.mu.Lock()
	for _, path := range paths {
		delete(a.files, path)
	}
	a.mu.Unlock()

	return nil
}

func (a *FakeAdapter) GetVisibility(ctx context.Context, path string) (entities.Visibility, error) {
	if a.GetVisibilityError != nil {
		return "", a.GetVisibilityError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return "", ErrFakeFileNotFound
	}

	return f.Visibility, nil
}

func (a *FakeAdapter) SetVisibility(ctx context.Context, path string, visibility entities.Visibility) error {
	if a.SetVisibilityError != nil {
		return a.SetVisibilityError
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	f, ok := a.files[path]
	if !ok {
		return ErrFakeFileNotFound
	}

	f.Visibility = visibility
	return nil
}

func (a *FakeAdapter) Files(ctx context.Context, directory string) ([]string, error) {
	if a.FilesError != nil {
		return nil, a.FilesError
	}
	return a.listFiles(directory, false), nil
}

func (a *FakeAdapter) AllFiles(ctx context.Context, directory string) ([]string, error) {
	if a.FilesError != nil {
		return nil, a.FilesError
	}
	return a.listFiles(directory, true), nil
}

func (a *FakeAdapter) listFiles(directory string, recursive bool) []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var files []string
	prefix := directory
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	for path := range a.files {
		if len(path) > len(prefix) && path[:len(prefix)] == prefix {
			rest := path[len(prefix):]
			hasSlash := false
			for _, c := range rest {
				if c == '/' {
					hasSlash = true
					break
				}
			}
			if recursive || !hasSlash {
				files = append(files, path)
			}
		}
	}

	return files
}

func (a *FakeAdapter) Directories(ctx context.Context, directory string) ([]string, error) {
	if a.DirectoriesError != nil {
		return nil, a.DirectoriesError
	}
	return a.listDirectories(directory, false), nil
}

func (a *FakeAdapter) AllDirectories(ctx context.Context, directory string) ([]string, error) {
	if a.DirectoriesError != nil {
		return nil, a.DirectoriesError
	}
	return a.listDirectories(directory, true), nil
}

func (a *FakeAdapter) listDirectories(directory string, recursive bool) []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	dirs := make(map[string]bool)
	prefix := directory
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	for path := range a.files {
		if len(path) > len(prefix) && path[:len(prefix)] == prefix {
			rest := path[len(prefix):]
			for i, c := range rest {
				if c == '/' {
					dir := prefix + rest[:i]
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

func (a *FakeAdapter) MakeDirectory(ctx context.Context, path string) error {
	if a.MakeDirectoryError != nil {
		return a.MakeDirectoryError
	}
	return nil
}

func (a *FakeAdapter) DeleteDirectory(ctx context.Context, directory string) error {
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

func (a *FakeAdapter) URL(path string) string {
	if a.BaseURL == "" {
		return path
	}
	return a.BaseURL + "/" + path
}

func (a *FakeAdapter) TemporaryURL(ctx context.Context, path string, expiration time.Duration) (string, error) {
	if a.TemporaryURLError != nil {
		return "", a.TemporaryURLError
	}

	return a.URL(path) + "?expires=" + time.Now().Add(expiration).Format(time.RFC3339), nil
}

// Reset clears all stored files
func (a *FakeAdapter) Reset() {
	a.mu.Lock()
	a.files = make(map[string]*fakeFile)
	a.mu.Unlock()
}

// FileCount returns the number of stored files
func (a *FakeAdapter) FileCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.files)
}

// HasFile checks if a specific file exists
func (a *FakeAdapter) HasFile(path string) bool {
	a.mu.RLock()
	_, ok := a.files[path]
	a.mu.RUnlock()
	return ok
}

// GetFileContent retrieves file content directly
func (a *FakeAdapter) GetFileContent(path string) ([]byte, bool) {
	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return f.Content, true
}

// ErrFakeFileNotFound is returned when a file does not exist
var ErrFakeFileNotFound = errors.New("file not found")
