package amazon_s3

import (
	"context"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	gomime "github.com/cubewise-code/go-mime"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/storage"
)

func (adapter Adapter) PutFile(ctx context.Context, input entities.StorageInput) (*entities.StorageObject, error) {
	extension := input.File.Extension()
	mimetype := gomime.TypeByExtension(extension)
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
