package fake

import (
	"bytes"
	"testing"
)

// AssertStored asserts that a file was stored at the given path
func (a *Adapter) AssertStored(t testing.TB, path string) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	if _, ok := a.files[path]; !ok {
		t.Errorf("Expected file to be stored at %q, but it was not", path)
	}
}

// AssertNotStored asserts that no file was stored at the given path
func (a *Adapter) AssertNotStored(t testing.TB, path string) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	if _, ok := a.files[path]; ok {
		t.Errorf("Expected file NOT to be stored at %q, but it was", path)
	}
}

// AssertStoredContent asserts that a file was stored with the expected content
func (a *Adapter) AssertStoredContent(t testing.TB, path string, expected []byte) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	file, ok := a.files[path]
	if !ok {
		t.Errorf("Expected file to be stored at %q, but it was not", path)
		return
	}

	if !bytes.Equal(file.Content, expected) {
		t.Errorf("File content at %q does not match.\nGot: %q\nWant: %q", path, file.Content, expected)
	}
}

// AssertDeleted asserts that files were deleted
func (a *Adapter) AssertDeleted(t testing.TB, paths ...string) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, path := range paths {
		found := false
		for _, deletedPaths := range a.DeleteCalls {
			for _, deletedPath := range deletedPaths {
				if deletedPath == path {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			t.Errorf("Expected file at %q to be deleted, but it was not", path)
		}
	}
}

// AssertNothingStored asserts that no files were stored
func (a *Adapter) AssertNothingStored(t testing.TB) {
	t.Helper()
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.PutFileCalls)+len(a.PutCalls) > 0 {
		t.Errorf("Expected no files to be stored, but %d put operations occurred", len(a.PutFileCalls)+len(a.PutCalls))
	}
}
