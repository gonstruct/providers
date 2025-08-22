package amazon_ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

type Adapter struct {
	Region   string
	Host     string
	Port     int
	Username string
	Password string
}

func (adapter Adapter) NewClient(context context.Context) (*sesv2.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context,
		config.WithRegion(adapter.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(adapter.Username, adapter.Password, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return sesv2.NewFromConfig(cfg), nil
}
