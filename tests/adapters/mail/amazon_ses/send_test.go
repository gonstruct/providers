package amazon_ses_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"testing"

	"github.com/gonstruct/providers/adapters/mail/amazon_ses"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/entities/mailables"
)

// sent is the part of SES v2's SendEmail request these tests read.
type sent struct {
	ConfigurationSetName *string
	Content              struct {
		Simple struct {
			Body struct {
				Html *struct{ Data string }
				Text *struct{ Data string }
			}
			Headers []struct{ Name, Value string }
		}
	}
}

// fakeSES points the SES client at a server that keeps the request it gets.
func fakeSES(t *testing.T) *sent {
	t.Helper()

	request := &sent{}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, incoming *http.Request) {
		body, _ := io.ReadAll(incoming.Body)
		if err := json.Unmarshal(body, request); err != nil {
			t.Errorf("SES got a body that is not JSON: %s", body)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"MessageId":"test"}`))
	}))
	t.Cleanup(server.Close)
	t.Setenv("AWS_ENDPOINT_URL_SESV2", server.URL)

	return request
}

func send(t *testing.T, text string, envelope mailables.Envelope) {
	t.Helper()

	envelope.From = mailables.Address("news@example.com", "Example")
	envelope.To = []*mail.Address{{Address: "user@example.com"}}
	envelope.Subject = "What shipped"

	var html, plain bytes.Buffer
	html.WriteString("<p>Hello</p>")
	plain.WriteString(text)

	adapter := amazon_ses.Adapter{Region: "eu-west-1", Username: "key", Password: "secret"}
	if err := adapter.Send(context.Background(), entities.MailInput{Envelope: envelope, Html: html, Text: plain}); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestSend_CarriesHeadersTextAndConfigurationSet(t *testing.T) {
	request := fakeSES(t)

	send(t, "Hello", mailables.Envelope{
		Headers: []mailables.Header{
			{Name: "List-Unsubscribe", Value: "<https://example.com/u/1>"},
			{Name: "List-Unsubscribe-Post", Value: "List-Unsubscribe=One-Click"},
		},
		ConfigurationSet: "marketing",
	})

	simple := request.Content.Simple
	if simple.Body.Text == nil || simple.Body.Text.Data != "Hello" {
		t.Fatalf("Text = %#v, want Hello", simple.Body.Text)
	}

	if simple.Body.Html == nil || simple.Body.Html.Data != "<p>Hello</p>" {
		t.Fatalf("Html = %#v, want <p>Hello</p>", simple.Body.Html)
	}

	if len(simple.Headers) != 2 || simple.Headers[0].Name != "List-Unsubscribe" ||
		simple.Headers[1].Value != "List-Unsubscribe=One-Click" {
		t.Fatalf("Headers = %#v, want both List-Unsubscribe headers in order", simple.Headers)
	}

	if request.ConfigurationSetName == nil || *request.ConfigurationSetName != "marketing" {
		t.Fatalf("ConfigurationSetName = %v, want marketing", request.ConfigurationSetName)
	}
}

func TestSend_LeavesOutWhatTheMailHasNot(t *testing.T) {
	request := fakeSES(t)

	send(t, "", mailables.Envelope{})

	if request.Content.Simple.Body.Text != nil {
		t.Fatalf("Text = %#v, want none", request.Content.Simple.Body.Text)
	}

	if request.Content.Simple.Headers != nil || request.ConfigurationSetName != nil {
		t.Fatalf("want no headers and no configuration set, got %#v and %v",
			request.Content.Simple.Headers, request.ConfigurationSetName)
	}
}
