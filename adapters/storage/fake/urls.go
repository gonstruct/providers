package fake

import (
	"context"
	"time"
)

func (a *Adapter) URL(path string) string {
	if a.BaseURL == "" {
		return path
	}
	return a.BaseURL + "/" + path
}

func (a *Adapter) TemporaryURL(ctx context.Context, path string, expiration time.Duration) (string, error) {
	if a.TemporaryURLError != nil {
		return "", a.TemporaryURLError
	}

	return a.URL(path) + "?expires=" + time.Now().Add(expiration).Format(time.RFC3339), nil
}
