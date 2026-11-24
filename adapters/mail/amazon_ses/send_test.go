package amazon_ses

import (
	"bytes"
	"net/mail"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/entities/mailables"
)

func input(text string, envelope mailables.Envelope) entities.MailInput {
	envelope.From = mailables.Address("news@example.com", "Example")
	envelope.To = []*mail.Address{{Address: "user@example.com"}}
	envelope.Subject = "What shipped"

	var html, plain bytes.Buffer
	html.WriteString("<p>Hello</p>")
	plain.WriteString(text)

	return entities.MailInput{Envelope: envelope, Html: html, Text: plain}
}

func TestBuildMessage_CarriesHeadersTextAndConfigurationSet(t *testing.T) {
	message, err := buildMessage(input("Hello", mailables.Envelope{
		Headers: []mailables.Header{
			{Name: "List-Unsubscribe", Value: "<https://example.com/u/1>"},
			{Name: "List-Unsubscribe-Post", Value: "List-Unsubscribe=One-Click"},
		},
		ConfigurationSet: "marketing",
	}))
	if err != nil {
		t.Fatalf("buildMessage() error = %v", err)
	}

	simple := message.Content.Simple
	if got := aws.ToString(simple.Body.Text.Data); got != "Hello" {
		t.Fatalf("Text = %q, want Hello", got)
	}

	if got := aws.ToString(simple.Body.Html.Data); got != "<p>Hello</p>" {
		t.Fatalf("Html = %q, want <p>Hello</p>", got)
	}

	if len(simple.Headers) != 2 ||
		aws.ToString(simple.Headers[0].Name) != "List-Unsubscribe" ||
		aws.ToString(simple.Headers[1].Value) != "List-Unsubscribe=One-Click" {
		t.Fatalf("Headers = %#v, want both List-Unsubscribe headers in order", simple.Headers)
	}

	if got := aws.ToString(message.ConfigurationSetName); got != "marketing" {
		t.Fatalf("ConfigurationSetName = %q, want marketing", got)
	}
}

func TestBuildMessage_LeavesOutWhatTheMailHasNot(t *testing.T) {
	message, err := buildMessage(input("", mailables.Envelope{}))
	if err != nil {
		t.Fatalf("buildMessage() error = %v", err)
	}

	if message.Content.Simple.Body.Text != nil {
		t.Fatalf("Text = %#v, want none", message.Content.Simple.Body.Text)
	}

	if message.Content.Simple.Headers != nil || message.ConfigurationSetName != nil {
		t.Fatalf("want no headers and no configuration set, got %#v and %v",
			message.Content.Simple.Headers, message.ConfigurationSetName)
	}
}
