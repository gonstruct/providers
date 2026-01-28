package fake

import (
	"context"

	"github.com/gonstruct/providers/entities"
)

func (a *Adapter) GetVisibility(ctx context.Context, path string) (entities.Visibility, error) {
	if a.GetVisibilityError != nil {
		return "", a.GetVisibilityError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return "", ErrFileNotFound
	}

	return f.Visibility, nil
}

func (a *Adapter) SetVisibility(ctx context.Context, path string, visibility entities.Visibility) error {
	if a.SetVisibilityError != nil {
		return a.SetVisibilityError
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	f, ok := a.files[path]
	if !ok {
		return ErrFileNotFound
	}

	f.Visibility = visibility

	return nil
}
