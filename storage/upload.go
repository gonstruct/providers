package storage

import (
	"io"
	"time"

	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/entities/file"
)

// Writing files

func PutFile(path string, file file.File, optionSlice ...Option) (*entities.StorageObject, error) {
	options := apply(optionSlice...)

	return options.Adapter.PutFile(options.Context, entities.StorageInput{
		ID:   options.GenerateUniqueID(),
		File: file,
		Path: path,
	})
}

func Put(path string, contents []byte, optionSlice ...Option) error {
	options := apply(optionSlice...)

	return options.Adapter.Put(options.Context, path, contents)
}

func PutStream(path string, stream io.Reader, optionSlice ...Option) error {
	options := apply(optionSlice...)

	return options.Adapter.PutStream(options.Context, path, stream)
}

// Reading files

func Get(path string, optionSlice ...Option) ([]byte, error) {
	options := apply(optionSlice...)

	return options.Adapter.Get(options.Context, path)
}

func GetStream(path string, optionSlice ...Option) (io.ReadCloser, error) {
	options := apply(optionSlice...)

	return options.Adapter.GetStream(options.Context, path)
}

func Exists(path string, optionSlice ...Option) (bool, error) {
	options := apply(optionSlice...)

	return options.Adapter.Exists(options.Context, path)
}

func Missing(path string, optionSlice ...Option) (bool, error) {
	options := apply(optionSlice...)

	return options.Adapter.Missing(options.Context, path)
}

// File metadata

func Size(path string, optionSlice ...Option) (int64, error) {
	options := apply(optionSlice...)

	return options.Adapter.Size(options.Context, path)
}

func LastModified(path string, optionSlice ...Option) (time.Time, error) {
	options := apply(optionSlice...)

	return options.Adapter.LastModified(options.Context, path)
}

func MimeType(path string, optionSlice ...Option) (string, error) {
	options := apply(optionSlice...)

	return options.Adapter.MimeType(options.Context, path)
}

// File operations

func Copy(from, to string, optionSlice ...Option) (*entities.StorageObject, error) {
	options := apply(optionSlice...)

	return options.Adapter.Copy(options.Context, from, to)
}

func Move(from, to string, optionSlice ...Option) (*entities.StorageObject, error) {
	options := apply(optionSlice...)

	return options.Adapter.Move(options.Context, from, to)
}

func Delete(paths []string, optionSlice ...Option) error {
	options := apply(optionSlice...)

	return options.Adapter.Delete(options.Context, paths...)
}

// Visibility

func GetVisibility(path string, optionSlice ...Option) (entities.Visibility, error) {
	options := apply(optionSlice...)

	return options.Adapter.GetVisibility(options.Context, path)
}

func SetVisibility(path string, visibility entities.Visibility, optionSlice ...Option) error {
	options := apply(optionSlice...)

	return options.Adapter.SetVisibility(options.Context, path, visibility)
}

// Directories

func Files(directory string, optionSlice ...Option) ([]string, error) {
	options := apply(optionSlice...)

	return options.Adapter.Files(options.Context, directory)
}

func AllFiles(directory string, optionSlice ...Option) ([]string, error) {
	options := apply(optionSlice...)

	return options.Adapter.AllFiles(options.Context, directory)
}

func Directories(directory string, optionSlice ...Option) ([]string, error) {
	options := apply(optionSlice...)

	return options.Adapter.Directories(options.Context, directory)
}

func AllDirectories(directory string, optionSlice ...Option) ([]string, error) {
	options := apply(optionSlice...)

	return options.Adapter.AllDirectories(options.Context, directory)
}

func MakeDirectory(path string, optionSlice ...Option) error {
	options := apply(optionSlice...)

	return options.Adapter.MakeDirectory(options.Context, path)
}

func DeleteDirectory(directory string, optionSlice ...Option) error {
	options := apply(optionSlice...)

	return options.Adapter.DeleteDirectory(options.Context, directory)
}

// URLs

func URL(path string, optionSlice ...Option) string {
	options := apply(optionSlice...)

	return options.Adapter.URL(path)
}

func TemporaryURL(path string, expiration time.Duration, optionSlice ...Option) (string, error) {
	options := apply(optionSlice...)

	return options.Adapter.TemporaryURL(options.Context, path, expiration)
}
