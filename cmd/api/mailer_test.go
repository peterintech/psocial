package main

import (
	"strings"
	"testing"

	"github.com/peterintech/psocial/internal/mailer"
)

func TestResolveNewMailer(t *testing.T) {
	t.Run("SMTP", func(t *testing.T) {
		client, err := resolveNewMailer(mailConfig{
			provider:  "smtp",
			fromName:  "Psocial",
			fromEmail: "sender@example.com",
			smtp: smtpConfig{
				host:     "smtp.gmail.com",
				port:     587,
				username: "sender@example.com",
				password: "app-password",
			},
		})
		if err != nil {
			t.Fatalf("resolveNewMailer returned an error: %v", err)
		}
		if _, ok := client.(*mailer.SMTPMailer); !ok {
			t.Fatalf("expected SMTPMailer, got %T", client)
		}
	})

	t.Run("SendGrid", func(t *testing.T) {
		client, err := resolveNewMailer(mailConfig{
			provider:  "sendgrid",
			fromName:  "Psocial",
			fromEmail: "sender@example.com",
			sendGrid: sendGridConfig{
				apiKey: "test-key",
			},
		})
		if err != nil {
			t.Fatalf("resolveNewMailer returned an error: %v", err)
		}
		if _, ok := client.(*mailer.SendGridMailer); !ok {
			t.Fatalf("expected SendGridMailer, got %T", client)
		}
	})

	t.Run("unsupported provider", func(t *testing.T) {
		_, err := resolveNewMailer(mailConfig{provider: "unknown"})
		if err == nil || !strings.Contains(err.Error(), "unsupported MAIL_PROVIDER") {
			t.Fatalf("expected unsupported provider error, got %v", err)
		}
	})
}
