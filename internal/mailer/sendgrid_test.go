package mailer

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type fakeSendGridClient struct {
	responses []*rest.Response
	errors    []error
	calls     int
}

func (c *fakeSendGridClient) Send(*mail.SGMailV3) (*rest.Response, error) {
	index := c.calls
	c.calls++

	var response *rest.Response
	if index < len(c.responses) {
		response = c.responses[index]
	}
	var err error
	if index < len(c.errors) {
		err = c.errors[index]
	}
	return response, err
}

func newTestSendGridMailer(t *testing.T, client sendGridClient) *SendGridMailer {
	t.Helper()

	m, err := NewSendGridMailer("Psocial", "sender@example.com", "test-api-key")
	if err != nil {
		t.Fatal(err)
	}
	m.client = client
	m.wait = func(context.Context, time.Duration) error { return nil }
	return m
}

func TestSendGridMailerRetriesNon2xxResponses(t *testing.T) {
	client := &fakeSendGridClient{responses: []*rest.Response{
		{StatusCode: 401},
		{StatusCode: 202},
	}}
	m := newTestSendGridMailer(t, client)

	if err := m.Send(context.Background(), UserWelcomeTemplate, "Peter", "recipient@example.com", welcomeData()); err != nil {
		t.Fatalf("expected retry to succeed: %v", err)
	}
	if client.calls != 2 {
		t.Fatalf("expected two SendGrid attempts, got %d", client.calls)
	}
}

func TestSendGridMailerFailsAfterNon2xxResponses(t *testing.T) {
	client := &fakeSendGridClient{responses: []*rest.Response{
		{StatusCode: 400},
		{StatusCode: 400},
		{StatusCode: 400},
	}}
	m := newTestSendGridMailer(t, client)

	err := m.Send(context.Background(), UserWelcomeTemplate, "Peter", "recipient@example.com", welcomeData())
	if err == nil || !strings.Contains(err.Error(), "status 400") {
		t.Fatalf("expected non-2xx delivery failure, got %v", err)
	}
	if client.calls != maxRetries {
		t.Fatalf("expected %d SendGrid attempts, got %d", maxRetries, client.calls)
	}
}

func TestSendGridMailerRetriesTransportErrors(t *testing.T) {
	client := &fakeSendGridClient{
		responses: []*rest.Response{nil, {StatusCode: 202}},
		errors:    []error{errors.New("network error"), nil},
	}
	m := newTestSendGridMailer(t, client)

	if err := m.Send(context.Background(), UserWelcomeTemplate, "Peter", "recipient@example.com", welcomeData()); err != nil {
		t.Fatalf("expected retry to succeed: %v", err)
	}
	if client.calls != 2 {
		t.Fatalf("expected two SendGrid attempts, got %d", client.calls)
	}
}

func TestNewSendGridMailerRejectsUnsafeFromName(t *testing.T) {
	_, err := NewSendGridMailer("Psocial\r\nBcc: attacker@example.com", "sender@example.com", "test-api-key")
	if err == nil || !strings.Contains(err.Error(), "MAIL_FROM_NAME") {
		t.Fatalf("expected unsafe sender name error, got %v", err)
	}
}
