package storage

import (
	"testing"

	"github.com/gonstruct/providers/entities/file"
	"github.com/gonstruct/providers/storage"
)

func TestPutFile(t *testing.T) {
	fake := storage.Fake()

	content := []byte("hello world")
	f := file.FromBytes("doc.pdf", content)

	obj, _ := storage.PutFile("uploads", f)

	fake.AssertStored(t, obj.Path)
	fake.AssertStoredContent(t, obj.Path, content)
}

func TestNothingStored(t *testing.T) {
	fake := storage.Fake()

	fake.AssertNothingStored(t)
}

func TestErrorInjection(t *testing.T) {
	fake := storage.Fake()
	fake.PutFileError = storage.Err("put file", nil)

	_, err := storage.PutFile("path", file.FromBytes("test.txt", []byte("x")))

	if err == nil {
		t.Error("expected error")
	}
}
