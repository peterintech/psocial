package mailer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendGridClient interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

type SendGridMailer struct {
	fromName  string
	fromEmail string
	apiKey    string
	client    sendGridClient
	wait      retryWaitFunc
}

func NewSendGridMailer(fromName, fromEmail, apiKey string) (*SendGridMailer, error) {
	fromName = strings.TrimSpace(fromName)
	fromEmail = strings.TrimSpace(fromEmail)
	apiKey = strings.TrimSpace(apiKey)

	if err := validateDisplayName(fromName); err != nil {
		return nil, fmt.Errorf("invalid MAIL_FROM_NAME: %w", err)
	}
	if err := validateAddress(fromEmail); err != nil {
		return nil, fmt.Errorf("invalid MAIL_FROM_EMAIL: %w", err)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("SENDGRID_API_KEY is required")
	}

	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromName:  fromName,
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
		wait:      waitForRetry,
	}, nil
}

func (s *SendGridMailer) Send(ctx context.Context, templateFile, username, email string, data any) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail(username, email)

	rendered, err := renderTemplate(templateFile, data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, rendered.subject, to, "", rendered.body)

	var retryErr error
	for i := range maxRetries {
		if err := ctx.Err(); err != nil {
			return err
		}

		response, err := s.client.Send(message)
		switch {
		case err != nil:
			retryErr = err
		case response == nil:
			retryErr = fmt.Errorf("SendGrid returned an empty response")
		case response.StatusCode < 200 || response.StatusCode >= 300:
			retryErr = fmt.Errorf("SendGrid returned status %d", response.StatusCode)
		default:
			return nil
		}

		if i < maxRetries-1 {
			if err := s.wait(ctx, time.Second*time.Duration(i+1)); err != nil {
				return err
			}
		}
	}
	return fmt.Errorf("failed to send email to %s after %d attempts: %w", email, maxRetries, retryErr)
}
