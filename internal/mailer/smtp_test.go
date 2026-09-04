package mailer

import (
	"context"
	"errors"
	"net/mail"
	"net/smtp"
	"strings"
	"testing"
	"time"
)

func validSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:      "smtp.gmail.com",
		Port:      587,
		Username:  "sender@example.com",
		Password:  "app-password",
		FromName:  "Psocial",
		FromEmail: "sender@example.com",
	}
}

func welcomeData() any {
	return struct {
		Username      string
		ActivationURL string
	}{Username: "Peter", ActivationURL: "https://example.com/activate?token=test"}
}

func TestNewSMTPMailerValidatesConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*SMTPConfig)
	}{
		{name: "host", mutate: func(cfg *SMTPConfig) { cfg.Host = "" }},
		{name: "port", mutate: func(cfg *SMTPConfig) { cfg.Port = 0 }},
		{name: "username", mutate: func(cfg *SMTPConfig) { cfg.Username = "" }},
		{name: "password", mutate: func(cfg *SMTPConfig) { cfg.Password = "" }},
		{name: "from name", mutate: func(cfg *SMTPConfig) { cfg.FromName = "" }},
		{name: "from name header injection", mutate: func(cfg *SMTPConfig) { cfg.FromName = "Psocial\r\nBcc: attacker@example.com" }},
		{name: "from email", mutate: func(cfg *SMTPConfig) { cfg.FromEmail = "invalid" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validSMTPConfig()
			tt.mutate(&cfg)
			if _, err := NewSMTPMailer(cfg); err == nil {
				t.Fatalf("expected invalid %s to return an error", tt.name)
			}
		})
	}
}

func TestSMTPMailerSendBuildsHTMLMessage(t *testing.T) {
	m, err := NewSMTPMailer(validSMTPConfig())
	if err != nil {
		t.Fatal(err)
	}

	var gotAddr, gotFrom string
	var gotTo []string
	var gotMessage string
	m.sendMail = func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
		if auth == nil {
			t.Fatal("expected SMTP authentication")
		}
		gotAddr, gotFrom, gotTo, gotMessage = addr, from, to, string(msg)
		return nil
	}

	if err := m.Send(context.Background(), UserWelcomeTemplate, "Peter", "recipient@example.com", welcomeData()); err != nil {
		t.Fatalf("Send returned an error: %v", err)
	}
	if gotAddr != "smtp.gmail.com:587" || gotFrom != "sender@example.com" {
		t.Fatalf("unexpected SMTP envelope: addr=%q from=%q", gotAddr, gotFrom)
	}
	if len(gotTo) != 1 || gotTo[0] != "recipient@example.com" {
		t.Fatalf("unexpected recipients: %#v", gotTo)
	}

	required := []string{
		"From: " + (&mail.Address{Name: "Psocial", Address: "sender@example.com"}).String() + "\r\n",
		"To: " + (&mail.Address{Name: "Peter", Address: "recipient@example.com"}).String() + "\r\n",
		"Subject: Finish Registration with Psocial\r\n",
		"MIME-Version: 1.0\r\n",
		"Content-Type: text/html; charset=\"UTF-8\"\r\n",
		"\r\n<!doctype html>",
	}
	for _, value := range required {
		if !strings.Contains(gotMessage, value) {
			t.Errorf("message did not contain %q", value)
		}
	}

	bodySeparator := strings.Index(gotMessage, "\r\n\r\n")
	if bodySeparator == -1 {
		t.Fatal("message did not contain an SMTP header/body separator")
	}
	body := gotMessage[bodySeparator+4:]
	if strings.Contains(strings.ReplaceAll(body, "\r\n", ""), "\n") {
		t.Error("message body contained a bare LF instead of SMTP CRLF line endings")
	}
}

func TestNormalizeSMTPBodyUsesCRLF(t *testing.T) {
	tests := map[string]string{
		"LF":    "\n<!doctype html>\n<body>hello</body>\n",
		"CRLF":  "\r\n<!doctype html>\r\n<body>hello</body>\r\n",
		"CR":    "\r<!doctype html>\r<body>hello</body>\r",
		"mixed": "\r\n<!doctype html>\n<body>hello</body>\r",
	}
	want := "<!doctype html>\r\n<body>hello</body>\r\n"

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if got := normalizeSMTPBody(input); got != want {
				t.Fatalf("normalizeSMTPBody() = %q, want %q", got, want)
			}
		})
	}
}

func TestSMTPMailerSendRetries(t *testing.T) {
	m, err := NewSMTPMailer(validSMTPConfig())
	if err != nil {
		t.Fatal(err)
	}

	attempts := 0
	m.sendMail = func(string, smtp.Auth, string, []string, []byte) error {
		attempts++
		return errors.New("SMTP unavailable")
	}
	m.wait = func(context.Context, time.Duration) error { return nil }

	err = m.Send(context.Background(), UserWelcomeTemplate, "Peter", "recipient@example.com", welcomeData())
	if err == nil {
		t.Fatal("expected delivery failure")
	}
	if attempts != maxRetries {
		t.Fatalf("expected %d attempts, got %d", maxRetries, attempts)
	}
}

func TestSMTPMailerSendHonorsCancelledContext(t *testing.T) {
	m, err := NewSMTPMailer(validSMTPConfig())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := m.Send(ctx, UserWelcomeTemplate, "Peter", "recipient@example.com", welcomeData()); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}

func TestSMTPMailerRejectsHeaderInjection(t *testing.T) {
	m, err := NewSMTPMailer(validSMTPConfig())
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Send(context.Background(), UserWelcomeTemplate, "Peter\r\nBcc: attacker@example.com", "recipient@example.com", welcomeData()); err == nil {
		t.Fatal("expected recipient header injection to be rejected")
	}
}
