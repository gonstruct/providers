package amazon_s3

import (
	"context"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	gomime "github.com/cubewise-code/go-mime"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/storage"
)

func (adapter Adapter) Copy(ctx context.Context, from, to string) (*entities.StorageObject, error) {
	name := path.Base(to)
	mimetype := gomime.TypeByExtension(path.Ext(to))

	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	if _, err := client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(adapter.Bucket),
		CopySource: aws.String(adapter.Bucket + "/" + from),
		Key:        aws.String(to),
	}); err != nil {
		return nil, storage.PathErr("copy", to, err)
	}

	return &entities.StorageObject{
		Name:     name,
		Path:     to,
		MimeType: mimetype,
	}, nil
}

func (adapter Adapter) Move(ctx context.Context, from, to string) (*entities.StorageObject, error) {
	name := path.Base(to)
	mimetype := gomime.TypeByExtension(path.Ext(to))

	if _, err := adapter.Copy(ctx, from, to); err != nil {
		return nil, err
	}

	if err := adapter.Delete(ctx, from); err != nil {
		return nil, err
	}

	return &entities.StorageObject{
		Name:     name,
		Path:     to,
		MimeType: mimetype,
	}, nil
}

func (adapter Adapter) Delete(ctx context.Context, paths ...string) error {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return storage.Err("create S3 client", err)
	}

	objects := make([]types.ObjectIdentifier, len(paths))
	for i, p := range paths {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(p),
		}
	}

	if _, err := client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(adapter.Bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	}); err != nil {
		return storage.Err("delete", err)
	}

	return nil
}
