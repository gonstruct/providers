package mail_test

import (
	"embed"
	"net/mail"
	"testing"

	"github.com/gonstruct/providers/adapters/mail/fake"
	"github.com/gonstruct/providers/entities/mailables"
	pmail "github.com/gonstruct/providers/mail"
)

//go:embed mail
var testTemplatesFS embed.FS

// testMailable implements contracts.Mailable for testing.
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

func TestSend_WithFakeAdapter(t *testing.T) {
	f := pmail.Fake(
		pmail.WithFakeTemplates(testTemplatesFS),
		pmail.WithFakeDefaultEnvelope(mailables.Envelope{
			From: mailables.Address("noreply@test.com", "Test App"),
		}),
	)

	mailable := testMailable{
		envelope: mailables.Envelope{
			Subject: "Test Email",
			To:      []*mail.Address{{Address: "test@example.com", Name: "Test User"}},
		},
		content: mailables.Content{
			View: "welcome.html",
			With: map[string]any{"name": "John"},
		},
	}

	err := pmail.Send(mailable)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	// Use assertion methods (Laravel-style)
	f.AssertSentCount(t, 1)
	f.AssertSentTo(t, "test@example.com")
	f.AssertSentWithSubject(t, "Test Email")
}

func TestSend_FakeRecordsCalls(t *testing.T) {
	f := pmail.Fake()

	// Verify nothing sent initially
	f.AssertNothingSent(t)

	if f.LastCall() != nil {
		t.Error("LastCall() should be nil before any calls")
	}

	// Verify Reset works
	f.Reset()

	if f.SentCount() != 0 {
		t.Error("SentCount() should be 0 after Reset()")
	}
}

func TestSend_AssertSentWithPredicate(t *testing.T) {
	f := pmail.Fake(
		pmail.WithFakeTemplates(testTemplatesFS),
		pmail.WithFakeDefaultEnvelope(mailables.Envelope{
			From: mailables.Address("noreply@test.com", "Test App"),
		}),
	)

	mailable := testMailable{
		envelope: mailables.Envelope{
			Subject: "Welcome!",
			To:      []*mail.Address{{Address: "user@example.com", Name: "User"}},
		},
		content: mailables.Content{
			View: "welcome.html",
			With: map[string]any{"name": "Alice"},
		},
	}

	err := pmail.Send(mailable)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	// Use custom predicate for complex assertions
	assertSentWith(t, f, func(call fake.SendCall) bool {
		return call.Input.Envelope.Subject == "Welcome!" &&
			len(call.Input.Envelope.To) > 0 &&
			call.Input.Envelope.To[0].Address == "user@example.com"
	})
}

// assertSentWith is a helper for custom predicate assertions.
func assertSentWith(t testing.TB, f *fake.Adapter, predicate func(call fake.SendCall) bool) {
	t.Helper()

	for _, call := range f.Calls {
		if predicate(call) {
			return
		}
	}

	t.Error("No email matched the predicate")
}

func TestEnvelope_Merge(t *testing.T) {
	base := &mailables.Envelope{
		Subject: "Original Subject",
		To:      mailables.Addresses("original@example.com"),
	}

	override := mailables.Envelope{
		Subject: "New Subject",
		To:      mailables.Addresses("additional@example.com"),
	}

	merged := base.Merge(override)

	if merged.Subject != "New Subject" {
		t.Errorf("Subject = %q, want %q", merged.Subject, "New Subject")
	}

	// To addresses should be combined
	if len(merged.To) != 2 {
		t.Errorf("To count = %d, want 2", len(merged.To))
	}
}
