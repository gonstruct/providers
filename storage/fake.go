package storage

import (
	"github.com/gonstruct/providers/adapters/storage/fake"
)

// Fake sets up a fake storage adapter for testing and returns it for assertions.
// This replaces any existing storage provider.
//
// Example:
//
//	func TestUpload(t *testing.T) {
//	    fake := storage.Fake()
//
//	    // Your code that uses storage.Put(), storage.Get(), etc.
//	    storage.Put(ctx, "path/to/file.txt", []byte("content"))
//
//	    // Assert
//	    fake.AssertStored(t, "path/to/file.txt")
//	}
func Fake() *fake.Adapter {
	adapter := fake.New()

	globalProvider = &provider{
		adapter: adapter,
	}

	return adapter
}
