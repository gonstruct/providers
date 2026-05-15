package mail_test

import (
	"context"
	"embed"
	"net/mail"
	"strings"
	"testing"

	"github.com/gonstruct/providers/contracts"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/entities/mailables"
	pmail "github.com/gonstruct/providers/mail"
)

//go:embed mail/*.html
var testTemplatesFS embed.FS

type testMailable struct {
	envelope    mailables.Envelope
	content     mailables.Content
	attachments mailables.AttachmentSlice
}

func (m testMailable) Envelope() mailables.Envelope {
	return m.envelope
}

func (m testMailable) Content() mailables.Content {
	return m.content
}

func (m testMailable) Attachments() mailables.AttachmentSlice {
	return m.attachments
}

type capturingAdapter struct {
	calls int
	ctx   context.Context
	input entities.MailInput
	err   error
}

func (a *capturingAdapter) Send(ctx context.Context, input entities.MailInput) error {
	a.calls++
	a.ctx = ctx
	a.input = input

	return a.err
}

func TestSend_MergesDefaultEnvelopeAndRendersTemplate(t *testing.T) {
	fake := pmail.Fake(
		pmail.WithFakeTemplates(testTemplatesFS),
		pmail.WithFakeDefaultEnvelope(mailables.Envelope{
			From:    mailables.Address("noreply@example.com", "Example App"),
			ReplyTo: mailables.Address("support@example.com", "Support"),
			Cc:      []*mail.Address{{Address: "default-cc@example.com"}},
			Bcc:     []*mail.Address{{Address: "default-bcc@example.com"}},
		}),
	)

	err := pmail.Send(testMailable{
		envelope: mailables.Envelope{
			Subject: "Welcome",
			To:      []*mail.Address{{Address: "user@example.com", Name: "User"}},
			Cc:      []*mail.Address{{Address: "message-cc@example.com"}},
		},
		content: mailables.Content{
			View: "welcome.html",
			With: map[string]any{"name": "User"},
		},
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	fake.AssertSentCount(t, 1)

	call := fake.Calls[0]
	if call.Input.Envelope.From == nil || call.Input.Envelope.From.Address != "noreply@example.com" {
		t.Fatalf("From = %#v, want noreply@example.com", call.Input.Envelope.From)
	}

	if call.Input.Envelope.ReplyTo == nil || call.Input.Envelope.ReplyTo.Address != "support@example.com" {
		t.Fatalf("ReplyTo = %#v, want support@example.com", call.Input.Envelope.ReplyTo)
	}

	if len(call.Input.Envelope.To) != 1 || call.Input.Envelope.To[0].Address != "user@example.com" {
		t.Fatalf("To = %#v, want user@example.com", call.Input.Envelope.To)
	}

	wantCc := []string{"default-cc@example.com", "message-cc@example.com"}
	if got := emailAddresses(call.Input.Envelope.Cc); !equalStrings(got, wantCc) {
		t.Fatalf("Cc = %v, want [default-cc@example.com message-cc@example.com]", got)
	}

	if got := emailAddresses(call.Input.Envelope.Bcc); !equalStrings(got, []string{"default-bcc@example.com"}) {
		t.Fatalf("Bcc = %v, want [default-bcc@example.com]", got)
	}

	if got := call.HTML; !strings.Contains(got, "Hello User") {
		t.Fatalf("HTML = %q, want rendered template content", got)
	}

	if call.Subject != "Welcome" {
		t.Fatalf("Subject = %q, want Welcome", call.Subject)
	}
}

func TestSend_DefaultEnvelopeRecipientsDoNotLeakBetweenSends(t *testing.T) {
	fake := pmail.Fake(
		pmail.WithFakeTemplates(testTemplatesFS),
		pmail.WithFakeDefaultEnvelope(mailables.Envelope{
			To:  []*mail.Address{{Address: "default-to@example.com"}},
			Cc:  []*mail.Address{{Address: "default-cc@example.com"}},
			Bcc: []*mail.Address{{Address: "default-bcc@example.com"}},
		}),
	)

	first := testMailable{
		envelope: mailables.Envelope{
			Subject: "First",
			To:      []*mail.Address{{Address: "first-to@example.com"}},
			Cc:      []*mail.Address{{Address: "first-cc@example.com"}},
			Bcc:     []*mail.Address{{Address: "first-bcc@example.com"}},
		},
		content: mailables.Content{View: "welcome.html", With: map[string]any{"name": "First"}},
	}

	second := testMailable{
		envelope: mailables.Envelope{
			Subject: "Second",
			To:      []*mail.Address{{Address: "second-to@example.com"}},
			Cc:      []*mail.Address{{Address: "second-cc@example.com"}},
			Bcc:     []*mail.Address{{Address: "second-bcc@example.com"}},
		},
		content: mailables.Content{View: "welcome.html", With: map[string]any{"name": "Second"}},
	}

	if err := pmail.Send(first); err != nil {
		t.Fatalf("first Send() error = %v", err)
	}

	if err := pmail.Send(second); err != nil {
		t.Fatalf("second Send() error = %v", err)
	}

	fake.AssertSentCount(t, 2)

	assertRecipients(t, fake.Calls[0].Input.Envelope.To, []string{"default-to@example.com", "first-to@example.com"})
	assertRecipients(t, fake.Calls[0].Input.Envelope.Cc, []string{"default-cc@example.com", "first-cc@example.com"})
	assertRecipients(t, fake.Calls[0].Input.Envelope.Bcc, []string{"default-bcc@example.com", "first-bcc@example.com"})

	assertRecipients(t, fake.Calls[1].Input.Envelope.To, []string{"default-to@example.com", "second-to@example.com"})
	assertRecipients(t, fake.Calls[1].Input.Envelope.Cc, []string{"default-cc@example.com", "second-cc@example.com"})
	assertRecipients(t, fake.Calls[1].Input.Envelope.Bcc, []string{"default-bcc@example.com", "second-bcc@example.com"})
}

func TestSend_UsesPerCallAdapterAndContext(t *testing.T) {
	providerAdapter := pmail.Fake(
		pmail.WithFakeTemplates(testTemplatesFS),
	)

	ctxKey := struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey, "trace-id")
	adapter := &capturingAdapter{}

	err := pmail.Send(
		testMailable{
			envelope: mailables.Envelope{
				Subject: "Overridden Adapter",
				To:      []*mail.Address{{Address: "context@example.com"}},
			},
			content: mailables.Content{View: "welcome.html", With: map[string]any{"name": "Context"}},
		},
		pmail.WithAdapter(adapter),
		pmail.WithContext(ctx),
	)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if adapter.calls != 1 {
		t.Fatalf("adapter calls = %d, want 1", adapter.calls)
	}

	if providerAdapter.SentCount() != 0 {
		t.Fatalf("provider adapter sent count = %d, want 0", providerAdapter.SentCount())
	}

	if got := adapter.ctx.Value(ctxKey); got != "trace-id" {
		t.Fatalf("context value = %v, want trace-id", got)
	}

	if got := adapter.input.Envelope.To[0].Address; got != "context@example.com" {
		t.Fatalf("To = %q, want context@example.com", got)
	}

	if got := adapter.input.Html.String(); !strings.Contains(got, "Hello Context") {
		t.Fatalf("HTML = %q, want rendered content", got)
	}
}

func TestSend_PassesAttachmentsThrough(t *testing.T) {
	fake := pmail.Fake(
		pmail.WithFakeTemplates(testTemplatesFS),
	)

	attachments := mailables.Attachments(
		mailables.Attachment(
			mailables.WithName("report.txt"),
			mailables.WithMime("text/plain"),
			mailables.WithContent([]byte("hello")),
		),
	)

	err := pmail.Send(testMailable{
		envelope: mailables.Envelope{
			Subject: "Attachments",
			To:      []*mail.Address{{Address: "files@example.com"}},
		},
		content:     mailables.Content{View: "welcome.html", With: map[string]any{"name": "Files"}},
		attachments: attachments,
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	call := fake.LastCall()
	if call == nil {
		t.Fatal("LastCall() = nil, want recorded call")
	}

	if len(call.Input.Attachments) != 1 {
		t.Fatalf("attachments count = %d, want 1", len(call.Input.Attachments))
	}

	if got := call.Input.Attachments[0].Name; got != "report.txt" {
		t.Fatalf("attachment name = %q, want report.txt", got)
	}

	if got := call.Input.Attachments[0].Mime; got != "text/plain" {
		t.Fatalf("attachment mime = %q, want text/plain", got)
	}

	if got := string(call.Input.Attachments[0].Content()); got != "hello" {
		t.Fatalf("attachment content = %q, want hello", got)
	}
}

func TestSend_ReturnsTemplateParseError(t *testing.T) {
	pmail.Fake()

	err := pmail.Send(testMailable{
		envelope: mailables.Envelope{
			Subject: "Missing Template",
			To:      []*mail.Address{{Address: "missing@example.com"}},
		},
		content: mailables.Content{View: "missing.html"},
	})
	if err == nil {
		t.Fatal("Send() error = nil, want template parse error")
	}

	if !strings.Contains(err.Error(), "failed to parse email template") {
		t.Fatalf("error = %q, want parse error", err)
	}
}

func assertRecipients(t testing.TB, addresses []*mail.Address, want []string) {
	t.Helper()

	got := emailAddresses(addresses)
	if !equalStrings(got, want) {
		t.Fatalf("addresses = %v, want %v", got, want)
	}
}

func emailAddresses(addresses []*mail.Address) []string {
	result := make([]string, len(addresses))
	for i, address := range addresses {
		result[i] = address.Address
	}

	return result
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}

	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}

	return true
}

var _ contracts.Mail = (*capturingAdapter)(nil)
