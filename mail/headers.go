package mail

import (
	"fmt"
	"strings"

	"github.com/gonstruct/providers/entities/mailables"
)

// reserved are the headers the envelope or the message itself writes. One
// set as an extra header would replace who the mail goes to, or break it.
var reserved = map[string]bool{
	"from": true, "to": true, "cc": true, "bcc": true, "reply-to": true, "sender": true,
	"subject": true, "date": true, "message-id": true, "return-path": true,
	"mime-version": true, "content-type": true, "content-transfer-encoding": true,
	"x-ses-configuration-set": true,
}

// CheckHeaders refuses headers an adapter must not send: a reserved name,
// a name that is not 1 to 126 printable ASCII characters without a colon
// (as SES takes them), and a value with a line break, which would start a
// header of its own.
func CheckHeaders(headers []mailables.Header) error {
	for _, header := range headers {
		if reserved[strings.ToLower(header.Name)] {
			return fmt.Errorf("%w: %s is set by the envelope", ErrHeader, header.Name)
		}

		if len(header.Name) == 0 || len(header.Name) > 126 {
			return fmt.Errorf("%w: a name of %d characters", ErrHeader, len(header.Name))
		}

		for _, character := range header.Name {
			if character < 33 || character > 126 || character == ':' {
				return fmt.Errorf("%w: %q is not a header name", ErrHeader, header.Name)
			}
		}

		if strings.ContainsAny(header.Value, "\r\n") {
			return fmt.Errorf("%w: %s has a line break", ErrHeader, header.Name)
		}
	}

	return nil
}
