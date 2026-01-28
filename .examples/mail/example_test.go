package mail

import (
	"testing"

	"github.com/gonstruct/providers/entities/mailables"
	pmail "github.com/gonstruct/providers/mail"
)

func TestSendEmail(t *testing.T) {
	fake := pmail.Fake(
		pmail.WithFakeDefaultEnvelope(mailables.Envelope{
			From: mailables.Address("noreply@app.com", "App"),
		}),
	)

	_ = pmail.Send(WelcomeMail{To: "user@example.com", Name: "John"})

	fake.AssertSentCount(t, 1)
	fake.AssertSentTo(t, "user@example.com")
	fake.AssertSentWithSubject(t, "Welcome!")
}

func TestNothingSent(t *testing.T) {
	fake := pmail.Fake()

	fake.AssertNothingSent(t)
}
