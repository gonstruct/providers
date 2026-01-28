package fake

import (
	"errors"
	"sync"
	"time"

	"github.com/gonstruct/providers/entities"
)

// ErrFileNotFound is returned when a file does not exist.
var ErrFileNotFound = errors.New("file not found")

// Adapter is an in-memory storage adapter for testing.
type Adapter struct {
	mu    sync.RWMutex
	files map[string]*fakeFile

	// Call tracking
	PutFileCalls []PutFileCall
	PutCalls     []PutCall
	GetCalls     []string
	DeleteCalls  [][]string
	CopyCalls    []CopyCall
	MoveCalls    []MoveCall

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

	// BaseURL for URL generation
	BaseURL string
}

type fakeFile struct {
	Content      []byte
	MimeType     string
	Visibility   entities.Visibility
	LastModified time.Time
}

type PutFileCall struct {
	Path    string
	Content []byte
}

type PutCall struct {
	Path    string
	Content []byte
}

type CopyCall struct {
	From string
	To   string
}

type MoveCall struct {
	From string
	To   string
}

// New creates a new fake storage adapter.
func New() *Adapter {
	return &Adapter{
		files:   make(map[string]*fakeFile),
		BaseURL: "https://fake.storage.test",
	}
}

// Reset clears all stored files and recorded calls.
func (a *Adapter) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.files = make(map[string]*fakeFile)
	a.PutFileCalls = nil
	a.PutCalls = nil
	a.GetCalls = nil
	a.DeleteCalls = nil
	a.CopyCalls = nil
	a.MoveCalls = nil
}

// --- Helper Methods ---

// HasFile returns true if a file exists at the given path.
func (a *Adapter) HasFile(path string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	_, ok := a.files[path]

	return ok
}

// GetFileContent returns the content of a file at the given path.
func (a *Adapter) GetFileContent(path string) ([]byte, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	file, ok := a.files[path]
	if !ok {
		return nil, false
	}

	return file.Content, true
}

// FileCount returns the number of stored files.
func (a *Adapter) FileCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return len(a.files)
}

// StoredPaths returns all stored file paths.
func (a *Adapter) StoredPaths() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	paths := make([]string, 0, len(a.files))
	for path := range a.files {
		paths = append(paths, path)
	}

	return paths
}
