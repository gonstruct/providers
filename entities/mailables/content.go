package mailables

import (
	"bytes"
	"embed"
	"fmt"
	"path"
	"text/template"
)

type Content struct {
	View string
	With map[string]any
}

func (content Content) Parse(templates embed.FS) (bytes.Buffer, error) {
	template, err := template.ParseFS(templates, path.Join("mail", content.View))
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err = template.Execute(&body, content.With); err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to execute email template: %w", err)
	}

	return body, nil
}
