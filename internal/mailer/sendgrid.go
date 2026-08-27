package mailer

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGridMailer(fromEmail, apiKey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (s *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, s.fromEmail)
	to := mail.NewEmail(username, email)

	//template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return 500, fmt.Errorf("failed to parse template: %w", err)
	}
	var subject, body bytes.Buffer

	err = tmpl.ExecuteTemplate(&subject, "subject", data)
	if err != nil {
		return 500, err
	}

	err = tmpl.ExecuteTemplate(&body, "body", data)
	if err != nil {
		return 500, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{SandboxMode: &mail.Setting{
		Enable: &isSandbox,
	}})

	var retryErr error
	for i := range maxRetries {
		res, retryErr := s.client.Send(message)
		if retryErr != nil {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return res.StatusCode, nil
	}
	return -1, fmt.Errorf("failed to send email to %s after %d attempts, error: %v", email, maxRetries, retryErr)
}
