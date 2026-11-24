package amazon_ses

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/mail"
)

func (adapter Adapter) Send(ctx context.Context, input entities.MailInput) error {
	message, err := buildMessage(input)
	if err != nil {
		return err
	}

	client, err := adapter.NewClient(ctx)
	if err != nil {
		return mail.Err("create SES client", err)
	}

	if _, err := client.SendEmail(ctx, message); err != nil {
		return mail.Err("send via SES", err)
	}

	return nil
}

// buildMessage is the input as SES takes it.
//
//nolint:cyclop,funlen
func buildMessage(input entities.MailInput) (*sesv2.SendEmailInput, error) {
	message := &sesv2.SendEmailInput{}

	var subject *types.Content
	if input.Envelope.Subject != "" {
		subject = &types.Content{
			Data:    aws.String(input.Envelope.Subject),
			Charset: aws.String("UTF-8"),
		}
	} else {
		return nil, mail.Err("validate", mail.ErrNoSubject)
	}

	if input.Envelope.From != nil {
		message.FromEmailAddress = aws.String(input.Envelope.From.String())
	} else {
		return nil, mail.Err("validate", mail.ErrNoSender)
	}

	message.Destination = &types.Destination{}
	if len(input.Envelope.To) > 0 {
		message.Destination.ToAddresses = input.Envelope.To.String()
	} else {
		return nil, mail.Err("validate", mail.ErrNoRecipients)
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

	body := &types.Body{
		Html: &types.Content{
			Data:    aws.String(input.Html.String()),
			Charset: aws.String("UTF-8"),
		},
	}
	if input.Text.Len() > 0 {
		body.Text = &types.Content{
			Data:    aws.String(input.Text.String()),
			Charset: aws.String("UTF-8"),
		}
	}

	headers := make([]types.MessageHeader, 0, len(input.Envelope.Headers))
	for _, header := range input.Envelope.Headers {
		headers = append(headers, types.MessageHeader{
			Name:  aws.String(header.Name),
			Value: aws.String(header.Value),
		})
	}

	message.Content = &types.EmailContent{
		Simple: &types.Message{
			Subject:     subject,
			Body:        body,
			Attachments: attachments,
		},
	}

	if len(headers) > 0 {
		message.Content.Simple.Headers = headers
	}

	if input.Envelope.ConfigurationSet != "" {
		message.ConfigurationSetName = aws.String(input.Envelope.ConfigurationSet)
	}

	return message, nil
}
