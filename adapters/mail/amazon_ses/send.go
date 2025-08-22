package amazon_ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/gonstruct/providers/entities"
)

//nolint:cyclop,funlen
func (adapter Adapter) Send(context context.Context, input entities.MailInput) error {
	message := &sesv2.SendEmailInput{}

	var subject *types.Content
	if input.Envelope.Subject != "" {
		subject = &types.Content{
			Data:    aws.String(input.Envelope.Subject),
			Charset: aws.String("UTF-8"),
		}
	} else {
		return fmt.Errorf("no subject specified")
	}

	if input.Envelope.From != nil {
		message.FromEmailAddress = aws.String(input.Envelope.From.String())
	} else {
		return fmt.Errorf("no sender specified")
	}

	message.Destination = &types.Destination{}
	if len(input.Envelope.To) > 0 {
		message.Destination.ToAddresses = input.Envelope.To.String()
	} else {
		return fmt.Errorf("no recipient(s) specified")
	}

	if input.Envelope.ReplyTo != nil {
		message.ReplyToAddresses = []string{input.Envelope.ReplyTo.String()}
	}

	if input.Envelope.Cc != nil {
		message.Destination.CcAddresses = input.Envelope.Cc.String()
	}

	if input.Envelope.Bcc != nil {
		message.Destination.BccAddresses = input.Envelope.Bcc.String()
	}

	var attachments []types.Attachment

	if len(input.Attachments) > 0 {
		for _, attachment := range input.Attachments {
			attachments = append(attachments, types.Attachment{
				FileName:                aws.String(attachment.Name),
				ContentType:             aws.String(attachment.Mime),
				ContentTransferEncoding: types.AttachmentContentTransferEncodingBase64,
				RawContent:              attachment.Content(),
			})
		}
	}

	message.Content = &types.EmailContent{
		Simple: &types.Message{
			Subject: subject,
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(input.Html.String()),
					Charset: aws.String("UTF-8"),
				},
				// Text: // TODOd: Add text/plain support if needed,
			},
			Attachments: attachments,
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
