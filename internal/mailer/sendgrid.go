package mailer

import (
	"context"
	"fmt"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromName  string
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGridMailer(fromName, fromEmail, apiKey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromName:  fromName,
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
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

		_, retryErr = s.client.Send(message)
		if retryErr != nil {
			if i < maxRetries-1 {
				timer := time.NewTimer(time.Second * time.Duration(i+1))
				select {
				case <-ctx.Done():
					timer.Stop()
					return ctx.Err()
				case <-timer.C:
				}
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to send email to %s after %d attempts: %w", email, maxRetries, retryErr)
}
