package amazon_s3

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	gomime "github.com/cubewise-code/go-mime"
	"github.com/gonstruct/providers/storage"
)

// Get retrieves the contents of a file.
func (adapter Adapter) Get(ctx context.Context, path string) ([]byte, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, storage.PathErr("get", path, err)
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

// GetStream returns a reader for the file contents.
func (adapter Adapter) GetStream(ctx context.Context, path string) (io.ReadCloser, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, storage.PathErr("get stream", path, err)
	}

	return result.Body, nil
}

// Exists checks if a file exists.
func (adapter Adapter) Exists(ctx context.Context, path string) (bool, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return false, storage.Err("create S3 client", err)
	}

	_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		// Check if it's a not found error
		return false, nil
	}

	return true, nil
}

// Missing checks if a file does not exist.
func (adapter Adapter) Missing(ctx context.Context, path string) (bool, error) {
	exists, err := adapter.Exists(ctx, path)

	return !exists, err
}

// Size returns the size of a file in bytes.
func (adapter Adapter) Size(ctx context.Context, path string) (int64, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return 0, storage.Err("create S3 client", err)
	}

	result, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return 0, storage.PathErr("get size", path, err)
	}

	if result.ContentLength != nil {
		return *result.ContentLength, nil
	}

	return 0, nil
}

// LastModified returns the last modification time of a file.
func (adapter Adapter) LastModified(ctx context.Context, path string) (time.Time, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return time.Time{}, storage.Err("create S3 client", err)
	}

	result, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return time.Time{}, storage.PathErr("get last modified", path, err)
	}

	if result.LastModified != nil {
		return *result.LastModified, nil
	}

	return time.Time{}, nil
}

// MimeType returns the MIME type of a file.
func (adapter Adapter) MimeType(ctx context.Context, path string) (string, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return "", storage.Err("create S3 client", err)
	}

	result, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return "", storage.PathErr("get mime type", path, err)
	}

	if result.ContentType != nil {
		return *result.ContentType, nil
	}

	// Fallback to extension-based detection
	return gomime.TypeByExtension(path), nil
}

// Put stores raw bytes at the given path.
func (adapter Adapter) Put(ctx context.Context, path string, contents []byte) error {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return storage.Err("create S3 client", err)
	}

	mimetype := gomime.TypeByExtension(path)
	if mimetype == "" {
		mimetype = "application/octet-stream"
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(adapter.Bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(contents),
		ContentType: aws.String(mimetype),
	})
	if err != nil {
		return storage.PathErr("put", path, err)
	}

	return nil
}

// PutStream stores content from a reader at the given path.
func (adapter Adapter) PutStream(ctx context.Context, path string, stream io.Reader) error {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return storage.Err("create S3 client", err)
	}

	mimetype := gomime.TypeByExtension(path)
	if mimetype == "" {
		mimetype = "application/octet-stream"
	}

	// Read stream into bytes (S3 SDK requires seekable body for retries)
	content, err := io.ReadAll(stream)
	if err != nil {
		return storage.Err("read stream", err)
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(adapter.Bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(mimetype),
	})
	if err != nil {
		return storage.PathErr("put stream", path, err)
	}

	return nil
}
