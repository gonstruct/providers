package fake

import (
	"testing"
)

// AssertSent asserts that at least one email was sent.
func (a *Adapter) AssertSent(t testing.TB) {
	t.Helper()

	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.Calls) == 0 {
		t.Error("Expected at least one email to be sent, but none were sent")
	}
}

// AssertSentCount asserts the exact number of emails sent.
func (a *Adapter) AssertSentCount(t testing.TB, count int) {
	t.Helper()

	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.Calls) != count {
		t.Errorf("Expected %d emails to be sent, got %d", count, len(a.Calls))
	}
}

// AssertSentTo asserts that an email was sent to the given recipient.
func (a *Adapter) AssertSentTo(t testing.TB, email string) {
	t.Helper()

	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, call := range a.Calls {
		for _, to := range call.To {
			if to == email {
				return
			}
		}
	}

	t.Errorf("Expected email to be sent to %q, but it was not", email)
}

// AssertNotSentTo asserts that no email was sent to the given recipient.
func (a *Adapter) AssertNotSentTo(t testing.TB, email string) {
	t.Helper()

	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, call := range a.Calls {
		for _, to := range call.To {
			if to == email {
				t.Errorf("Expected email NOT to be sent to %q, but it was", email)

				return
			}
		}
	}
}

// AssertSentWithSubject asserts that an email was sent with the given subject.
func (a *Adapter) AssertSentWithSubject(t testing.TB, subject string) {
	t.Helper()

	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, call := range a.Calls {
		if call.Subject == subject {
			return
		}
	}

	t.Errorf("Expected email to be sent with subject %q, but it was not", subject)
}

// AssertSentFrom asserts that an email was sent from the given address.
func (a *Adapter) AssertSentFrom(t testing.TB, email string) {
	t.Helper()

	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, call := range a.Calls {
		if call.From == email {
			return
		}
	}

	t.Errorf("Expected email to be sent from %q, but it was not", email)
}

// AssertNothingSent asserts that no emails were sent.
func (a *Adapter) AssertNothingSent(t testing.TB) {
	t.Helper()

	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.Calls) > 0 {
		t.Errorf("Expected no emails to be sent, but %d were sent", len(a.Calls))
	}
}
