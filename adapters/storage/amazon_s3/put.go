package s3

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	gomime "github.com/cubewise-code/go-mime"
	"github.com/gonstruct/providers/entities"
)

func (adapter Adapter) PutFile(context context.Context, input entities.StorageInput) (*entities.StorageObject, error) {
	extension := input.File.Extension()
	mimetype := gomime.TypeByExtension(extension)
	key := path.Join(input.Path, input.ID+extension)

	client, err := adapter.NewClient(context)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	if _, err := client.PutObject(context, &s3.PutObjectInput{
		Bucket:      aws.String(adapter.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(mimetype),
		Body:        input.File.Body,
	}); err != nil {
		return nil, err
	}

	return &entities.StorageObject{
		Name:     input.Name(),
		Path:     key,
		MimeType: mimetype,
	}, nil
}
