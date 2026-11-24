package mail_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/gonstruct/providers/entities/mailables"
	pmail "github.com/gonstruct/providers/mail"
)

func TestCheckHeaders(t *testing.T) {
	allowed := []mailables.Header{
		{Name: "List-Unsubscribe", Value: "<https://example.com/u/1>, <mailto:u@example.com>"},
		{Name: "List-Unsubscribe-Post", Value: "List-Unsubscribe=One-Click"},
	}
	if err := pmail.CheckHeaders(allowed); err != nil {
		t.Fatalf("CheckHeaders() = %v, want the unsubscribe headers allowed", err)
	}

	refused := map[string]mailables.Header{
		"a reserved name":         {Name: "bcc", Value: "someone@example.com"},
		"the configuration set":   {Name: "X-SES-Configuration-Set", Value: "other"},
		"an empty name":           {Name: "", Value: "x"},
		"a name with a colon":     {Name: "X-A:B", Value: "x"},
		"a name with a space":     {Name: "X A", Value: "x"},
		"a name with a line feed": {Name: "X-A\r\nBcc", Value: "x"},
		"a name too long":         {Name: strings.Repeat("X", 127), Value: "x"},
		"a value with a break":    {Name: "X-A", Value: "x\r\nBcc: someone@example.com"},
	}
	for name, header := range refused {
		if err := pmail.CheckHeaders([]mailables.Header{header}); !errors.Is(err, pmail.ErrHeader) {
			t.Errorf("%s: CheckHeaders() = %v, want ErrHeader", name, err)
		}
	}
}
