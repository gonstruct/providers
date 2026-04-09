package amazon_s3

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/storage"
)

// URL returns the public URL for a file.
func (adapter Adapter) URL(path string) string {
	if adapter.Endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", adapter.Endpoint, adapter.Bucket, path)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", adapter.Bucket, adapter.Region, path)
}

// TemporaryURL generates a presigned URL with an expiration time.
func (adapter Adapter) TemporaryURL(
	ctx context.Context,
	input entities.TemporaryStorageInput,
) (*entities.TemporaryStorageObject, error) {
	extension := input.File.Extension()

	key := path.Join(input.Path, input.ID+extension)

	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	result, err := s3.NewPresignClient(client).PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(input.Expiry))
	if err != nil {
		return nil, storage.PathErr("generate presigned url", key, err)
	}

	return &entities.TemporaryStorageObject{
		Name:         input.Name(),
		Path:         key,
		URL:          result.URL,
		Method:       result.Method,
		SignedHeader: result.SignedHeader,
	}, nil
}

func (adapter Adapter) TemporaryUploadURL(
	ctx context.Context,
	input entities.TemporaryStorageInput,
) (*entities.TemporaryStorageObject, error) {
	extension := input.File.Extension()
	mimetype := input.File.MimeType()

	key := path.Join(input.Path, input.ID+extension)

	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	result, err := s3.NewPresignClient(client).PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(adapter.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(mimetype),
		ACL:         storageVisibilityToS3ACL(input.Visibility),
	}, s3.WithPresignExpires(input.Expiry))
	if err != nil {
		return nil, storage.PathErr("generate presigned upload url", key, err)
	}

	return &entities.TemporaryStorageObject{
		Name:         input.Name(),
		Path:         key,
		URL:          result.URL,
		Method:       result.Method,
		SignedHeader: result.SignedHeader,
	}, nil
}
