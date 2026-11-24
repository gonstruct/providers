package smtp_test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"net"
	"net/mail"
	"strconv"
	"strings"
	"testing"

	"github.com/gonstruct/providers/adapters/mail/smtp"
	"github.com/gonstruct/providers/entities"
	"github.com/gonstruct/providers/entities/mailables"
	pmail "github.com/gonstruct/providers/mail"
)

// fakeSMTP is a server that takes one message and keeps it, with the
// addresses it was sent to.
type fakeSMTP struct {
	port       int
	recipients []string
	message    chan *mail.Message
}

func serve(t *testing.T) *fakeSMTP {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })

	server := &fakeSMTP{port: listener.Addr().(*net.TCPAddr).Port, message: make(chan *mail.Message, 1)}

	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}
		defer connection.Close()

		server.converse(t, connection)
	}()

	return server
}

func (server *fakeSMTP) converse(t *testing.T, connection net.Conn) {
	t.Helper()

	reader := bufio.NewReader(connection)
	reply := func(line string) { _, _ = connection.Write([]byte(line + "\r\n")) }

	reply("220 fake")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		command := strings.ToUpper(strings.TrimSpace(line))

		switch {
		case strings.HasPrefix(command, "EHLO"), strings.HasPrefix(command, "HELO"):
			reply("250 fake")
		case strings.HasPrefix(command, "RCPT TO:"):
			server.recipients = append(server.recipients, strings.Trim(strings.TrimSpace(line)[8:], "<>"))
			reply("250 ok")
		case command == "DATA":
			reply("354 go on")

			var data bytes.Buffer

			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil || dataLine == ".\r\n" {
					break
				}

				data.WriteString(dataLine)
			}

			message, err := mail.ReadMessage(&data)
			if err != nil {
				t.Errorf("not a message: %v", err)
			}

			server.message <- message

			reply("250 kept")
		case command == "QUIT":
			reply("221 bye")

			return
		default:
			reply("250 ok")
		}
	}
}

func send(t *testing.T, server *fakeSMTP, text string, envelope mailables.Envelope) error {
	t.Helper()

	envelope.From = mailables.Address("news@example.com", "Example")
	envelope.To = []*mail.Address{{Address: "user@example.com"}}
	envelope.Subject = "What shipped"

	var html, plain bytes.Buffer
	html.WriteString("<p>Hello</p>")
	plain.WriteString(text)

	adapter := &smtp.Adapter{Host: "127.0.0.1", Port: server.port}
	input := entities.MailInput{Envelope: envelope, Html: html, Text: plain}

	return adapter.Send(context.Background(), input)
}

func TestSend_CarriesHeadersTextAndConfigurationSet(t *testing.T) {
	server := serve(t)

	err := send(t, server, "Hello", mailables.Envelope{
		Headers: []mailables.Header{
			{Name: "List-Unsubscribe", Value: "<https://example.com/u/1>"},
			{Name: "List-Unsubscribe-Post", Value: "List-Unsubscribe=One-Click"},
		},
		ConfigurationSet: "marketing",
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	message := <-server.message
	for name, want := range map[string]string{
		"List-Unsubscribe":        "<https://example.com/u/1>",
		"List-Unsubscribe-Post":   "List-Unsubscribe=One-Click",
		"X-Ses-Configuration-Set": "marketing",
	} {
		if got := message.Header.Get(name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}

	contentType := message.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/alternative") {
		t.Fatalf("Content-Type = %q, want multipart/alternative", contentType)
	}

	var body bytes.Buffer
	_, _ = body.ReadFrom(message.Body)

	text, html := strings.Index(body.String(), "Hello"), strings.Index(body.String(), "<p>Hello</p>")
	if text < 0 || html < 0 || text > html {
		t.Fatalf("want the text part first and the HTML after it, got %s", strconv.Quote(body.String()))
	}
}

func TestSend_WithoutTextSendsHTMLAlone(t *testing.T) {
	server := serve(t)

	if err := send(t, server, "", mailables.Envelope{}); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	message := <-server.message
	if got := message.Header.Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", got)
	}

	if got := message.Header.Get("X-Ses-Configuration-Set"); got != "" {
		t.Fatalf("X-SES-CONFIGURATION-SET = %q, want none", got)
	}
}

func TestSend_RefusesAHeaderThatWouldAddARecipient(t *testing.T) {
	server := serve(t)

	err := send(t, server, "", mailables.Envelope{
		Headers: []mailables.Header{{Name: "Bcc", Value: "someone@example.com"}},
	})
	if !errors.Is(err, pmail.ErrHeader) {
		t.Fatalf("Send() error = %v, want ErrHeader", err)
	}

	if len(server.recipients) != 0 {
		t.Fatalf("recipients = %v, want none", server.recipients)
	}
}
