package amazon_s3

import (
	"context"
	"fmt"
	"time"

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
	path string,
	expiration time.Duration,
) (*entities.PresignedObject, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	presignClient := s3.NewPresignClient(client)

	result, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return nil, storage.PathErr("generate presigned url", path, err)
	}

	return &entities.PresignedObject{
		URL:          result.URL,
		Method:       result.Method,
		SignedHeader: result.SignedHeader,
	}, nil
}

func (adapter Adapter) TemporaryUploadURL(
	ctx context.Context,
	path string,
	expiration time.Duration,
) (*entities.PresignedObject, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	presignClient := s3.NewPresignClient(client)

	result, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(path),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return nil, storage.PathErr("generate presigned upload url", path, err)
	}

	return &entities.PresignedObject{
		URL:          result.URL,
		Method:       result.Method,
		SignedHeader: result.SignedHeader,
	}, nil
}
