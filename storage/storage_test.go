package storage_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/gonstruct/providers/entities/file"
	"github.com/gonstruct/providers/storage"
)

func TestPutFile(t *testing.T) {
	fake := storage.Fake()

	content := []byte("test file content")
	f := file.FromBytes("document.pdf", content)

	obj, err := storage.PutFile("uploads", f)
	if err != nil {
		t.Fatalf("PutFile() error = %v", err)
	}

	if obj.Name != "document" {
		t.Errorf("Name = %q, want %q", obj.Name, "document")
	}

	// Use assertion methods
	fake.AssertStored(t, obj.Path)

	stored, ok := fake.GetFileContent(obj.Path)
	if !ok {
		t.Fatal("could not get file content from fake")
	}

	if !bytes.Equal(stored, content) {
		t.Errorf("stored content = %q, want %q", stored, content)
	}
}

func TestPutFile_WithAdapter(t *testing.T) {
	// For inline adapter usage, we can still use the fake directly
	fake := storage.Fake()

	content := []byte("inline adapter test")
	f := file.FromBytes("test.txt", content)

	obj, err := storage.PutFile("path", f)
	if err != nil {
		t.Fatalf("PutFile() error = %v", err)
	}

	fake.AssertStored(t, obj.Path)
}

func TestPutFile_WithContext(t *testing.T) {
	storage.Fake()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	f := file.FromBytes("test.txt", []byte("content"))

	// The fake doesn't check context, but this tests the option works
	_, err := storage.PutFile("path", f, storage.WithContext(ctx))
	if err != nil {
		t.Fatalf("PutFile() error = %v", err)
	}
}

func TestPutFile_WithUniqueIDGenerator(t *testing.T) {
	fake := storage.Fake()

	f := file.FromBytes("test.txt", []byte("content"))

	obj, err := storage.PutFile("path", f,
		storage.WithUniqueIDGenerator(func() string { return "custom-id" }),
	)
	if err != nil {
		t.Fatalf("PutFile() error = %v", err)
	}

	// Path should contain our custom ID
	expected := "path/custom-id.txt"
	if obj.Path != expected {
		t.Errorf("Path = %q, want %q", obj.Path, expected)
	}

	fake.AssertStored(t, expected)
}

func TestPutFile_Error(t *testing.T) {
	fake := storage.Fake()
	fake.PutFileError = context.DeadlineExceeded

	f := file.FromBytes("test.txt", []byte("content"))

	_, err := storage.PutFile("path", f)
	if err == nil {
		t.Error("PutFile() should return error when adapter fails")
	}
}

func TestStorage_AssertNothingStored(t *testing.T) {
	fake := storage.Fake()

	// No operations performed
	fake.AssertNothingStored(t)
}

func TestStorage_AssertStoredContent(t *testing.T) {
	fake := storage.Fake()

	content := []byte("important data")
	f := file.FromBytes("data.bin", content)

	obj, err := storage.PutFile("files", f)
	if err != nil {
		t.Fatalf("PutFile() error = %v", err)
	}

	fake.AssertStoredContent(t, obj.Path, content)
}
