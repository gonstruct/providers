package fake

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/gonstruct/providers/entities"
)

func (a *Adapter) PutFile(ctx context.Context, input entities.StorageInput) (*entities.StorageObject, error) {
	if a.PutFileError != nil {
		return nil, a.PutFileError
	}

	content, err := io.ReadAll(input.File.Body)
	if err != nil {
		return nil, err
	}

	// Reset the reader
	if seeker, ok := input.File.Body.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	path := input.Path + "/" + input.ID + input.File.Extension()
	now := time.Now()

	a.mu.Lock()
	a.PutFileCalls = append(a.PutFileCalls, PutFileCall{Path: path, Content: content})
	a.files[path] = &fakeFile{
		Content:      content,
		MimeType:     "application/octet-stream",
		Visibility:   entities.VisibilityPrivate,
		LastModified: now,
	}
	a.mu.Unlock()

	return &entities.StorageObject{
		Name:         input.Name(),
		Path:         path,
		MimeType:     "application/octet-stream",
		Size:         int64(len(content)),
		LastModified: now,
		Visibility:   entities.VisibilityPrivate,
	}, nil
}

func (a *Adapter) Put(ctx context.Context, path string, contents []byte) error {
	if a.PutError != nil {
		return a.PutError
	}

	a.mu.Lock()
	a.PutCalls = append(a.PutCalls, PutCall{Path: path, Content: contents})
	a.files[path] = &fakeFile{
		Content:      contents,
		MimeType:     "application/octet-stream",
		Visibility:   entities.VisibilityPrivate,
		LastModified: time.Now(),
	}
	a.mu.Unlock()

	return nil
}

func (a *Adapter) PutStream(ctx context.Context, path string, stream io.Reader) error {
	if a.PutStreamError != nil {
		return a.PutStreamError
	}

	content, err := io.ReadAll(stream)
	if err != nil {
		return err
	}

	a.mu.Lock()
	a.files[path] = &fakeFile{
		Content:      content,
		MimeType:     "application/octet-stream",
		Visibility:   entities.VisibilityPrivate,
		LastModified: time.Now(),
	}
	a.mu.Unlock()

	return nil
}

func (a *Adapter) Get(ctx context.Context, path string) ([]byte, error) {
	if a.GetError != nil {
		return nil, a.GetError
	}

	a.mu.Lock()
	a.GetCalls = append(a.GetCalls, path)
	a.mu.Unlock()

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return nil, ErrFileNotFound
	}

	return f.Content, nil
}

func (a *Adapter) GetStream(ctx context.Context, path string) (io.ReadCloser, error) {
	if a.GetStreamError != nil {
		return nil, a.GetStreamError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return nil, ErrFileNotFound
	}

	return io.NopCloser(bytes.NewReader(f.Content)), nil
}

func (a *Adapter) Exists(ctx context.Context, path string) (bool, error) {
	if a.ExistsError != nil {
		return false, a.ExistsError
	}

	a.mu.RLock()
	_, ok := a.files[path]
	a.mu.RUnlock()

	return ok, nil
}

func (a *Adapter) Missing(ctx context.Context, path string) (bool, error) {
	exists, err := a.Exists(ctx, path)

	return !exists, err
}

func (a *Adapter) Size(ctx context.Context, path string) (int64, error) {
	if a.SizeError != nil {
		return 0, a.SizeError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return 0, ErrFileNotFound
	}

	return int64(len(f.Content)), nil
}

func (a *Adapter) LastModified(ctx context.Context, path string) (time.Time, error) {
	if a.LastModifiedError != nil {
		return time.Time{}, a.LastModifiedError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return time.Time{}, ErrFileNotFound
	}

	return f.LastModified, nil
}

func (a *Adapter) MimeType(ctx context.Context, path string) (string, error) {
	if a.MimeTypeError != nil {
		return "", a.MimeTypeError
	}

	a.mu.RLock()
	f, ok := a.files[path]
	a.mu.RUnlock()

	if !ok {
		return "", ErrFileNotFound
	}

	return f.MimeType, nil
}

func (a *Adapter) Copy(ctx context.Context, from, to string) error {
	if a.CopyError != nil {
		return a.CopyError
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.CopyCalls = append(a.CopyCalls, CopyCall{From: from, To: to})

	f, ok := a.files[from]
	if !ok {
		return ErrFileNotFound
	}

	a.files[to] = &fakeFile{
		Content:      append([]byte(nil), f.Content...),
		MimeType:     f.MimeType,
		Visibility:   f.Visibility,
		LastModified: time.Now(),
	}

	return nil
}

func (a *Adapter) Move(ctx context.Context, from, to string) error {
	if a.MoveError != nil {
		return a.MoveError
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.MoveCalls = append(a.MoveCalls, MoveCall{From: from, To: to})

	f, ok := a.files[from]
	if !ok {
		return ErrFileNotFound
	}

	a.files[to] = f
	delete(a.files, from)

	return nil
}

func (a *Adapter) Delete(ctx context.Context, paths ...string) error {
	if a.DeleteError != nil {
		return a.DeleteError
	}

	a.mu.Lock()

	a.DeleteCalls = append(a.DeleteCalls, paths)
	for _, path := range paths {
		delete(a.files, path)
	}

	a.mu.Unlock()

	return nil
}
