package local

import (
	"context"
	"errors"
	"time"

	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/storage"
)

// URL returns the public URL for a file
// For local storage, this requires a BaseURL to be configured.
func (a *Adapter) URL(filePath string) string {
	if a.BaseURL == "" {
		return ""
	}

	return a.BaseURL + "/" + filePath
}

// TemporaryURL generates a presigned URL with an expiration time.
// Local storage does not support presigned URLs natively, so this method returns an error.
func (a *Adapter) TemporaryURL(
	ctx context.Context,
	filePath string,
	expiration time.Duration,
) (*entities.PresignedObject, error) {
	return nil, storage.Err(
		"temporary url",
		errors.New("Local storage does not support temporary URLs; implement signed URL logic in your application"),
	)
}

// TemporaryUploadURL generates a presigned URL for uploading with an expiration time.
// Local storage does not support presigned URLs natively, so this method returns an error.
func (a *Adapter) TemporaryUploadURL(
	ctx context.Context,
	filePath string,
	expiration time.Duration,
) (*entities.PresignedObject, error) {
	return nil, storage.Err(
		"temporary upload url",
		errors.New("Local storage does not support temporary upload URLs; implement signed URL logic in your application"),
	)
}
