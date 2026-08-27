package mailer

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

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

func (s *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, s.fromEmail)
	to := mail.NewEmail(username, email)

	//template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	var subject, body bytes.Buffer

	err = tmpl.ExecuteTemplate(&subject, "subject", data)
	if err != nil {
		return err
	}

	err = tmpl.ExecuteTemplate(&body, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{SandboxMode: &mail.Setting{
		Enable: &isSandbox,
	}})

	for range maxRetries {
		res, err := s.client.Send(message)
		if err != nil {
			continue
		}
		log.Printf("Email sent to %s with status code: %d", email, res.StatusCode)
		return nil
	}
	return fmt.Errorf("failed to send email to %s after %d attempts", email, maxRetries)
}
