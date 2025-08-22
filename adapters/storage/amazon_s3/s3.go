package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Adapter struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	Endpoint        string
	UsePathStyle    bool
}

func (adapter Adapter) NewClient(context context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context,
		config.WithRegion(adapter.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(adapter.AccessKeyID, adapter.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = adapter.UsePathStyle
		o.BaseEndpoint = aws.String(adapter.Endpoint)
	}), nil
}
