package amazon_s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gonstruct/providers/storage"
)

// Copy copies a file from one location to another.
func (adapter Adapter) Copy(ctx context.Context, from, to string) error {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return storage.Err("create S3 client", err)
	}

	_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(adapter.Bucket),
		CopySource: aws.String(adapter.Bucket + "/" + from),
		Key:        aws.String(to),
	})
	if err != nil {
		return storage.PathErr("copy", from+" -> "+to, err)
	}

	return nil
}

// Move moves a file from one location to another.
func (adapter Adapter) Move(ctx context.Context, from, to string) error {
	// S3 doesn't have native move, so copy then delete
	if err := adapter.Copy(ctx, from, to); err != nil {
		return err
	}

	return adapter.Delete(ctx, from)
}

// Delete removes one or more files.
func (adapter Adapter) Delete(ctx context.Context, paths ...string) error {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return storage.Err("create S3 client", err)
	}

	// Use batch delete for efficiency
	objects := make([]types.ObjectIdentifier, len(paths))
	for i, path := range paths {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(path),
		}
	}

	_, err = client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(adapter.Bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		return storage.Err("delete objects", err)
	}

	return nil
}
