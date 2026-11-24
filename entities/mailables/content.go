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
	// TextView is the plain text version, rendered with the same With and
	// sent beside the HTML. Without one the mail is HTML alone.
	TextView string
}

func (content Content) Parse(templates embed.FS) (bytes.Buffer, error) {
	return content.parse(templates, content.View)
}

// ParseText renders the plain text version, nothing when there is none.
func (content Content) ParseText(templates embed.FS) (bytes.Buffer, error) {
	if content.TextView == "" {
		return bytes.Buffer{}, nil
	}

	return content.parse(templates, content.TextView)
}

func (content Content) parse(templates embed.FS, view string) (bytes.Buffer, error) {
	template, err := template.ParseFS(templates, path.Join("mail", view))
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err = template.Execute(&body, content.With); err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to execute email template: %w", err)
	}

	return body, nil
}
