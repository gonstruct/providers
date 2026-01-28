package local

import (
	"context"
	"errors"
	"path"
	"time"

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

// TemporaryURL generates a temporary URL for the file
// Local storage doesn't natively support signed URLs, so this returns an error
// You can implement signed URL support using your application's routing.
func (a *Adapter) TemporaryURL(ctx context.Context, filePath string, expiration time.Duration) (string, error) {
	// Check if file exists first
	exists, err := a.Exists(ctx, filePath)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", storage.PathErr("temporary url", filePath, storage.ErrFileNotFound)
	}

	// Local storage doesn't support temporary URLs natively
	// Applications can implement this via their routing layer
	if a.BaseURL == "" {
		return "", storage.Err(
			"temporary url",
			errors.New("BaseURL not configured; implement signed URL logic in your application"),
		)
	}

	// Return the regular URL - implement signing in your application
	return path.Join(a.BaseURL, filePath), nil
}
