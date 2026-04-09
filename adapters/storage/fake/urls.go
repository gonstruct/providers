package fake

import (
	"context"
	"fmt"
	"time"

	"github.com/gonstruct/providers/entities"
)

func (a *Adapter) URL(path string) string {
	if a.BaseURL == "" {
		return path
	}

	return a.BaseURL + "/" + path
}

func (a *Adapter) TemporaryURL(ctx context.Context, path string, expiration time.Duration) (*entities.PresignedObject, error) {
	if a.TemporaryURLError != nil {
		return nil, a.TemporaryURLError
	}

	return &entities.PresignedObject{
		URL:    a.URL(path) + "?expires=" + time.Now().Add(expiration).Format(time.RFC3339),
		Method: "GET",
	}, nil
}

func (a *Adapter) TemporaryUploadURL(ctx context.Context, path string, expiration time.Duration) (*entities.PresignedObject, error) {
	if a.TemporaryUploadURLError != nil {
		return nil, a.TemporaryUploadURLError
	}

	return nil, fmt.Errorf("Temporary upload URLs are not supported in the fake adapter")
}
