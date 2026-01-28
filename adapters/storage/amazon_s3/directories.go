package amazon_s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gonstruct/providers/storage"
)

// Files returns a list of files in the given directory (non-recursive)
func (adapter Adapter) Files(ctx context.Context, directory string) ([]string, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	prefix := directory
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	result, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(adapter.Bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, storage.PathErr("list files", directory, err)
	}

	var files []string
	for _, obj := range result.Contents {
		if obj.Key != nil && *obj.Key != prefix {
			files = append(files, *obj.Key)
		}
	}

	return files, nil
}

// AllFiles returns a list of all files in the directory and subdirectories
func (adapter Adapter) AllFiles(ctx context.Context, directory string) ([]string, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	prefix := directory
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	var files []string
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(adapter.Bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, storage.PathErr("list all files", directory, err)
		}

		for _, obj := range page.Contents {
			if obj.Key != nil && !strings.HasSuffix(*obj.Key, "/") {
				files = append(files, *obj.Key)
			}
		}
	}

	return files, nil
}

// Directories returns a list of directories in the given directory (non-recursive)
func (adapter Adapter) Directories(ctx context.Context, directory string) ([]string, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	prefix := directory
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	result, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(adapter.Bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, storage.PathErr("list directories", directory, err)
	}

	var dirs []string
	for _, prefix := range result.CommonPrefixes {
		if prefix.Prefix != nil {
			// Remove trailing slash for consistency
			dir := strings.TrimSuffix(*prefix.Prefix, "/")
			dirs = append(dirs, dir)
		}
	}

	return dirs, nil
}

// AllDirectories returns a list of all directories in the directory and subdirectories
func (adapter Adapter) AllDirectories(ctx context.Context, directory string) ([]string, error) {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return nil, storage.Err("create S3 client", err)
	}

	prefix := directory
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	dirSet := make(map[string]bool)
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(adapter.Bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, storage.PathErr("list all directories", directory, err)
		}

		for _, obj := range page.Contents {
			if obj.Key != nil {
				// Extract all parent directories
				parts := strings.Split(*obj.Key, "/")
				for i := 1; i < len(parts); i++ {
					dir := strings.Join(parts[:i], "/")
					if dir != "" && dir != directory {
						dirSet[dir] = true
					}
				}
			}
		}
	}

	dirs := make([]string, 0, len(dirSet))
	for dir := range dirSet {
		dirs = append(dirs, dir)
	}

	return dirs, nil
}

// MakeDirectory creates a directory (S3 doesn't have real directories, creates a placeholder)
func (adapter Adapter) MakeDirectory(ctx context.Context, path string) error {
	client, err := adapter.NewClient(ctx)
	if err != nil {
		return storage.Err("create S3 client", err)
	}

	// S3 doesn't have real directories, but we can create a placeholder
	key := path
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(adapter.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return storage.PathErr("create directory", path, err)
	}

	return nil
}

// DeleteDirectory removes a directory and all its contents
func (adapter Adapter) DeleteDirectory(ctx context.Context, directory string) error {
	// First, list all objects with the directory prefix
	files, err := adapter.AllFiles(ctx, directory)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	return adapter.Delete(ctx, files...)
}
