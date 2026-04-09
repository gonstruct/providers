package fake

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/gonstruct/providers/entities"
)

func (a *Adapter) URL(path string) string {
	if a.BaseURL == "" {
		return path
	}

	return a.BaseURL + "/" + path
}

func (a *Adapter) TemporaryURL(
	ctx context.Context,
	input entities.TemporaryStorageInput,
) (*entities.TemporaryStorageObject, error) {
	if a.TemporaryURLError != nil {
		return nil, a.TemporaryURLError
	}

	key := path.Join(input.Path, input.ID+input.File.Extension())

	return &entities.TemporaryStorageObject{
		Name:   input.Name(),
		Path:   key,
		URL:    a.URL(key) + "?expires=" + time.Now().Add(input.Expiry).Format(time.RFC3339),
		Method: "GET",
	}, nil
}

func (a *Adapter) TemporaryUploadURL(
	ctx context.Context,
	input entities.TemporaryStorageInput,
) (*entities.TemporaryStorageObject, error) {
	if a.TemporaryUploadURLError != nil {
		return nil, a.TemporaryUploadURLError
	}

	return nil, fmt.Errorf("Temporary upload URLs are not supported in the fake adapter")
}
