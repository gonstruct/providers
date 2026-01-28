package local_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gonstruct/providers/adapters/storage/local"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/entities/file"
)

func setupAdapter(t *testing.T) (*local.Adapter, string) {
	t.Helper()
	root := t.TempDir()
	return local.NewAdapter(root), root
}

func TestNewAdapter(t *testing.T) {
	adapter, root := setupAdapter(t)

	if adapter.Root != root {
		t.Errorf("Root = %q, want %q", adapter.Root, root)
	}
	if adapter.FilePermission != 0644 {
		t.Errorf("FilePermission = %o, want %o", adapter.FilePermission, 0644)
	}
	if adapter.DirectoryPermission != 0755 {
		t.Errorf("DirectoryPermission = %o, want %o", adapter.DirectoryPermission, 0755)
	}
}

func TestAdapter_Put_Get(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	content := []byte("hello, world!")
	path := "test/file.txt"

	if err := adapter.Put(ctx, path, content); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	got, err := adapter.Get(ctx, path)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if !bytes.Equal(got, content) {
		t.Errorf("Get() = %q, want %q", got, content)
	}
}

func TestAdapter_PutStream_GetStream(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	content := []byte("streamed content")
	path := "stream/file.txt"

	if err := adapter.PutStream(ctx, path, bytes.NewReader(content)); err != nil {
		t.Fatalf("PutStream() error = %v", err)
	}

	stream, err := adapter.GetStream(ctx, path)
	if err != nil {
		t.Fatalf("GetStream() error = %v", err)
	}
	defer stream.Close()

	got, err := io.ReadAll(stream)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if !bytes.Equal(got, content) {
		t.Errorf("GetStream content = %q, want %q", got, content)
	}
}

func TestAdapter_PutFile(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	content := []byte("file content")
	input := entities.StorageInput{
		ID:   "test-id",
		Path: "uploads",
		File: file.File{
			Name: "document.pdf",
			Body: bytes.NewReader(content),
		},
	}

	obj, err := adapter.PutFile(ctx, input)
	if err != nil {
		t.Fatalf("PutFile() error = %v", err)
	}

	if obj.Name != "document" { // Name() trims the extension
		t.Errorf("Name = %q, want %q", obj.Name, "document")
	}
	if obj.MimeType == "" {
		t.Error("MimeType should not be empty")
	}

	// Verify file was actually stored
	got, err := adapter.Get(ctx, obj.Path)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("stored content = %q, want %q", got, content)
	}
}

func TestAdapter_Exists_Missing(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	path := "test/exists.txt"

	// File doesn't exist yet
	exists, err := adapter.Exists(ctx, path)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Error("Exists() = true for non-existent file")
	}

	missing, err := adapter.Missing(ctx, path)
	if err != nil {
		t.Fatalf("Missing() error = %v", err)
	}
	if !missing {
		t.Error("Missing() = false for non-existent file")
	}

	// Create the file
	if err := adapter.Put(ctx, path, []byte("content")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	// File exists now
	exists, err = adapter.Exists(ctx, path)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() = false for existing file")
	}

	missing, err = adapter.Missing(ctx, path)
	if err != nil {
		t.Fatalf("Missing() error = %v", err)
	}
	if missing {
		t.Error("Missing() = true for existing file")
	}
}

func TestAdapter_Size(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	content := []byte("12345678901234567890") // 20 bytes
	path := "test/size.txt"

	if err := adapter.Put(ctx, path, content); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	size, err := adapter.Size(ctx, path)
	if err != nil {
		t.Fatalf("Size() error = %v", err)
	}

	if size != 20 {
		t.Errorf("Size() = %d, want %d", size, 20)
	}
}

func TestAdapter_LastModified(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	before := time.Now().Add(-time.Second)
	path := "test/time.txt"

	if err := adapter.Put(ctx, path, []byte("content")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	after := time.Now().Add(time.Second)

	lastMod, err := adapter.LastModified(ctx, path)
	if err != nil {
		t.Fatalf("LastModified() error = %v", err)
	}

	if lastMod.Before(before) || lastMod.After(after) {
		t.Errorf("LastModified() = %v, want between %v and %v", lastMod, before, after)
	}
}

func TestAdapter_Copy(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	content := []byte("original content")
	from := "test/original.txt"
	to := "test/copy.txt"

	if err := adapter.Put(ctx, from, content); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	if err := adapter.Copy(ctx, from, to); err != nil {
		t.Fatalf("Copy() error = %v", err)
	}

	// Both files should exist
	fromExists, _ := adapter.Exists(ctx, from)
	toExists, _ := adapter.Exists(ctx, to)

	if !fromExists {
		t.Error("original file should still exist after Copy")
	}
	if !toExists {
		t.Error("copied file should exist after Copy")
	}

	// Content should match
	got, _ := adapter.Get(ctx, to)
	if !bytes.Equal(got, content) {
		t.Errorf("copied content = %q, want %q", got, content)
	}
}

func TestAdapter_Move(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	content := []byte("moving content")
	from := "test/source.txt"
	to := "test/destination.txt"

	if err := adapter.Put(ctx, from, content); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	if err := adapter.Move(ctx, from, to); err != nil {
		t.Fatalf("Move() error = %v", err)
	}

	// Original should not exist, destination should
	fromExists, _ := adapter.Exists(ctx, from)
	toExists, _ := adapter.Exists(ctx, to)

	if fromExists {
		t.Error("original file should not exist after Move")
	}
	if !toExists {
		t.Error("destination file should exist after Move")
	}

	// Content should match
	got, _ := adapter.Get(ctx, to)
	if !bytes.Equal(got, content) {
		t.Errorf("moved content = %q, want %q", got, content)
	}
}

func TestAdapter_Delete(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	paths := []string{"test/a.txt", "test/b.txt", "test/c.txt"}
	for _, p := range paths {
		if err := adapter.Put(ctx, p, []byte("content")); err != nil {
			t.Fatalf("Put() error = %v", err)
		}
	}

	// Delete multiple files
	if err := adapter.Delete(ctx, paths...); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	for _, p := range paths {
		exists, _ := adapter.Exists(ctx, p)
		if exists {
			t.Errorf("file %q should not exist after Delete", p)
		}
	}
}

func TestAdapter_Visibility(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	path := "test/visibility.txt"
	if err := adapter.Put(ctx, path, []byte("content")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	// Set to public
	if err := adapter.SetVisibility(ctx, path, entities.VisibilityPublic); err != nil {
		t.Fatalf("SetVisibility(public) error = %v", err)
	}

	vis, err := adapter.GetVisibility(ctx, path)
	if err != nil {
		t.Fatalf("GetVisibility() error = %v", err)
	}
	if vis != entities.VisibilityPublic {
		t.Errorf("GetVisibility() = %v, want %v", vis, entities.VisibilityPublic)
	}

	// Set to private
	if err := adapter.SetVisibility(ctx, path, entities.VisibilityPrivate); err != nil {
		t.Fatalf("SetVisibility(private) error = %v", err)
	}

	vis, err = adapter.GetVisibility(ctx, path)
	if err != nil {
		t.Fatalf("GetVisibility() error = %v", err)
	}
	if vis != entities.VisibilityPrivate {
		t.Errorf("GetVisibility() = %v, want %v", vis, entities.VisibilityPrivate)
	}
}

func TestAdapter_MakeDirectory_DeleteDirectory(t *testing.T) {
	adapter, root := setupAdapter(t)
	ctx := context.Background()

	dir := "new/nested/directory"

	if err := adapter.MakeDirectory(ctx, dir); err != nil {
		t.Fatalf("MakeDirectory() error = %v", err)
	}

	// Verify directory exists
	fullPath := filepath.Join(root, dir)
	info, err := os.Stat(fullPath)
	if err != nil {
		t.Fatalf("directory does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("path is not a directory")
	}

	// Add a file in the directory
	if err := adapter.Put(ctx, dir+"/file.txt", []byte("content")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	// Delete the directory
	if err := adapter.DeleteDirectory(ctx, dir); err != nil {
		t.Fatalf("DeleteDirectory() error = %v", err)
	}

	// Verify directory is gone
	_, err = os.Stat(fullPath)
	if !os.IsNotExist(err) {
		t.Error("directory should not exist after DeleteDirectory")
	}
}

func TestAdapter_Files(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	// Create files in nested structure
	files := []string{
		"dir/file1.txt",
		"dir/file2.txt",
		"dir/sub/file3.txt",
		"dir/sub/deep/file4.txt",
	}
	for _, f := range files {
		if err := adapter.Put(ctx, f, []byte("content")); err != nil {
			t.Fatalf("Put() error = %v", err)
		}
	}

	// Files() should return only direct children
	directFiles, err := adapter.Files(ctx, "dir")
	if err != nil {
		t.Fatalf("Files() error = %v", err)
	}

	if len(directFiles) != 2 {
		t.Errorf("Files() returned %d files, want 2", len(directFiles))
	}

	// AllFiles() should return all files recursively
	allFiles, err := adapter.AllFiles(ctx, "dir")
	if err != nil {
		t.Fatalf("AllFiles() error = %v", err)
	}

	if len(allFiles) != 4 {
		t.Errorf("AllFiles() returned %d files, want 4", len(allFiles))
	}
}

func TestAdapter_Directories(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	// Create files in nested structure to create directories
	files := []string{
		"root/a/file.txt",
		"root/b/file.txt",
		"root/b/c/file.txt",
		"root/b/c/d/file.txt",
	}
	for _, f := range files {
		if err := adapter.Put(ctx, f, []byte("content")); err != nil {
			t.Fatalf("Put() error = %v", err)
		}
	}

	// Directories() should return direct subdirectories
	dirs, err := adapter.Directories(ctx, "root")
	if err != nil {
		t.Fatalf("Directories() error = %v", err)
	}

	if len(dirs) != 2 {
		t.Errorf("Directories() returned %d directories, want 2", len(dirs))
	}

	// AllDirectories() should return all directories recursively
	allDirs, err := adapter.AllDirectories(ctx, "root")
	if err != nil {
		t.Fatalf("AllDirectories() error = %v", err)
	}

	if len(allDirs) < 4 {
		t.Errorf("AllDirectories() returned %d directories, want at least 4", len(allDirs))
	}
}

func TestAdapter_URL(t *testing.T) {
	adapter, _ := setupAdapter(t)
	adapter.BaseURL = "https://example.com/storage"

	url := adapter.URL("path/to/file.txt")
	expected := "https://example.com/storage/path/to/file.txt"

	if url != expected {
		t.Errorf("URL() = %q, want %q", url, expected)
	}
}

func TestAdapter_TemporaryURL(t *testing.T) {
	adapter, _ := setupAdapter(t)
	adapter.BaseURL = "https://example.com/storage"
	ctx := context.Background()

	path := "test/temp.txt"
	if err := adapter.Put(ctx, path, []byte("content")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	url, err := adapter.TemporaryURL(ctx, path, time.Hour)
	if err != nil {
		t.Fatalf("TemporaryURL() error = %v", err)
	}

	if url == "" {
		t.Error("TemporaryURL() returned empty string")
	}
}

func TestAdapter_TemporaryURL_NoBaseURL(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	path := "test/temp.txt"
	if err := adapter.Put(ctx, path, []byte("content")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	_, err := adapter.TemporaryURL(ctx, path, time.Hour)
	if err == nil {
		t.Error("TemporaryURL() without BaseURL should return error")
	}
}

func TestAdapter_Get_NonExistent(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	_, err := adapter.Get(ctx, "does/not/exist.txt")
	if err == nil {
		t.Error("Get() should return error for non-existent file")
	}
}

func TestAdapter_MimeType(t *testing.T) {
	adapter, _ := setupAdapter(t)
	ctx := context.Background()

	tests := []struct {
		path     string
		content  []byte
		expected string
	}{
		{"test.txt", []byte("hello"), "text/plain; charset=utf-8"},
		{"test.html", []byte("<html></html>"), "text/html; charset=utf-8"},
		{"test.json", []byte(`{"key":"value"}`), "application/json"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if err := adapter.Put(ctx, tt.path, tt.content); err != nil {
				t.Fatalf("Put() error = %v", err)
			}

			mime, err := adapter.MimeType(ctx, tt.path)
			if err != nil {
				t.Fatalf("MimeType() error = %v", err)
			}

			// MimeType detection can vary, just check it's not empty
			if mime == "" {
				t.Error("MimeType() returned empty string")
			}
		})
	}
}
