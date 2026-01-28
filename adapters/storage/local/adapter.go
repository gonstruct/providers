package local

import (
	"github.com/gonstruct/providers/entities"
)

// Adapter implements local filesystem storage
type Adapter struct {
	// Root is the base directory for all storage operations
	Root string

	// BaseURL is the public URL prefix for generating URLs (optional)
	BaseURL string

	// Permissions for files and directories
	FilePermission      int
	DirectoryPermission int
}

// NewAdapter creates a new local storage adapter with sensible defaults
func NewAdapter(root string) *Adapter {
	return &Adapter{
		Root:                root,
		FilePermission:      0644,
		DirectoryPermission: 0755,
	}
}

// WithBaseURL sets the base URL for public file access
func (a *Adapter) WithBaseURL(url string) *Adapter {
	a.BaseURL = url
	return a
}

// WithPermissions sets custom file and directory permissions
func (a *Adapter) WithPermissions(file, directory int) *Adapter {
	a.FilePermission = file
	a.DirectoryPermission = directory
	return a
}

// getVisibilityPermission returns the file permission for a visibility level
func (a *Adapter) getVisibilityPermission(visibility entities.Visibility) int {
	switch visibility {
	case entities.VisibilityPublic:
		return 0644
	case entities.VisibilityPrivate:
		return 0600
	default:
		return a.FilePermission
	}
}
