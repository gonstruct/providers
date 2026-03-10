package amazon_s3

import (
	"bytes"
	"context"
	"io"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	gomime "github.com/cubewise-code/go-mime"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/storage"
)

func (adapter Adapter) PutFile(ctx context.Context, input entities.StorageInput) (*entities.StorageObject, error) {
	extension := input.File.Extension()
	mimetype := input.File.MimeType()

	key := path.Join(input.Path, input.ID+extension)

	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	if _, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(adapter.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(mimetype),
		Body:        input.File.Body,
	}); err != nil {
		return nil, storage.PathErr("put file", key, err)
	}

	return &entities.StorageObject{
		Name:     input.Name(),
		Path:     key,
		MimeType: mimetype,
	}, nil
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
