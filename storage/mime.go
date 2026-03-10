package storage

import (
	"path/filepath"
	"strings"

	gomime "github.com/cubewise-code/go-mime"
)

func NewMimeTypes(allowed ...string) MimeTypes {
	return MimeTypes{allowed: allowed}
}

type MimeTypes struct {
	allowed []string
}

func (m *MimeTypes) IsAccepted(mimeType string, def ...string) bool {
	if len(m.allowed) == 0 {
		return true
	}

	if mimeType == "" && len(def) > 0 {
		mimeType = def[0]
	}

	for _, allowed := range m.allowed {
		if allowed == "*" || allowed == mimeType {
			return true
		}

		if strings.HasSuffix(allowed, "/*") {
			prefix := strings.TrimSuffix(allowed, "/*") + "/"
			if strings.HasPrefix(mimeType, prefix) {
				return true
			}
		}
	}

	return false
}

func (m *MimeTypes) IsAcceptedPath(path string, def ...string) bool {
	mimetype := gomime.TypeByExtension(filepath.Ext(path))

	return m.IsAccepted(mimetype, def...)
}
