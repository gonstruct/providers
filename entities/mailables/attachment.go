package mailables

type attachment struct {
	Name    string
	Mime    string
	content []byte
}

func (a attachment) Content() []byte {
	if a.content == nil {
		return []byte{}
	}

	return a.content
}

type AttachmentSlice []attachment

func Attachments(attachments ...attachment) AttachmentSlice {
	return attachments
}

type attachmentOption func(*attachment)

func Attachment(options ...attachmentOption) attachment {
	attachment := attachment{}

	for _, option := range options {
		option(&attachment)
	}

	return attachment
}

func WithName(name string) attachmentOption {
	return func(a *attachment) {
		a.Name = name
	}
}

func WithMime(mime string) attachmentOption {
	return func(a *attachment) {
		a.Mime = mime
	}
}

func WithContent(content []byte) attachmentOption {
	return func(a *attachment) {
		a.content = content
	}
}
