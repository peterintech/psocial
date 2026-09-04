package mailer

import (
	"context"
	"os"
	"strconv"
	"testing"
)

func TestSMTPIntegration(t *testing.T) {
	if os.Getenv("SMTP_INTEGRATION_TEST") != "1" {
		t.Skip("set SMTP_INTEGRATION_TEST=1 to send a real email")
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		t.Fatalf("SMTP_PORT must be a number: %v", err)
	}
	recipient := os.Getenv("SMTP_TEST_RECIPIENT")
	if recipient == "" {
		t.Fatal("SMTP_TEST_RECIPIENT is required")
	}

	client, err := NewSMTPMailer(SMTPConfig{
		Host:      os.Getenv("SMTP_HOST"),
		Port:      port,
		Username:  os.Getenv("SMTP_USERNAME"),
		Password:  os.Getenv("SMTP_PASSWORD"),
		FromName:  os.Getenv("MAIL_FROM_NAME"),
		FromEmail: os.Getenv("MAIL_FROM_EMAIL"),
	})
	if err != nil {
		t.Fatalf("invalid SMTP integration configuration: %v", err)
	}

	data := struct {
		Username      string
		ActivationURL string
	}{
		Username:      "SMTP Test",
		ActivationURL: "https://example.invalid/psocial-smtp-test",
	}
	if err := client.Send(context.Background(), UserWelcomeTemplate, "SMTP Test", recipient, data); err != nil {
		t.Fatalf("SMTP delivery failed: %v", err)
	}
}
