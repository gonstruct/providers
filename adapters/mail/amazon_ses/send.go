package amazon_ses

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/gonstruct/providers/entities/mailables"
)

func (adapter Adapter) Send(context context.Context, envelope mailables.Envelope, body bytes.Buffer) error {
	message := &sesv2.SendEmailInput{}

	var subject *types.Content
	if envelope.Subject != "" {
		subject = &types.Content{
			Data:    aws.String(envelope.Subject),
			Charset: aws.String("UTF-8"),
		}
	} else {
		return fmt.Errorf("no subject specified")
	}

	if envelope.From != nil {
		message.FromEmailAddress = aws.String(envelope.From.String())
	} else {
		return fmt.Errorf("no sender specified")
	}

	message.Destination = &types.Destination{}
	if len(envelope.To) > 0 {
		message.Destination.ToAddresses = envelope.To.String()
	} else {
		return fmt.Errorf("no recipient(s) specified")
	}

	if envelope.ReplyTo != nil {
		message.ReplyToAddresses = []string{envelope.ReplyTo.String()}
	}

	if envelope.Cc != nil {
		message.Destination.CcAddresses = envelope.Cc.String()
	}

	if envelope.Bcc != nil {
		message.Destination.BccAddresses = envelope.Bcc.String()
	}

	// var attachments []types.Attachment
	// if options.Attachments != nil && len(options.Attachments) > 0 {
	// 	for _, attachment := range options.Attachments {
	// 		attachments = append(attachments, types.Attachment{
	// 			FileName:                aws.String(attachment.Name),
	// 			ContentType:             aws.String(attachment.Mime),
	// 			ContentTransferEncoding: types.AttachmentContentTransferEncodingBase64,
	// 			RawContent:              attachment.Content(),
	// 		})
	// 	}
	// }

	message.Content = &types.EmailContent{
		Simple: &types.Message{
			Subject: subject,
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(body.String()),
					Charset: aws.String("UTF-8"),
				},
				// Text: // TODOd: Add text/plain support if needed,
			},
			// Attachments: attachments,
		},
	}

	client, err := adapter.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create SES client: %w", err)
	}

	if _, err := client.SendEmail(context, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
