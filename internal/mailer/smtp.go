package mailer

import (
	"bytes"
	"context"
	"fmt"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromName  string
	FromEmail string
}

type smtpSendMailFunc func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error
type retryWaitFunc func(context.Context, time.Duration) error

type SMTPMailer struct {
	host      string
	addr      string
	username  string
	password  string
	fromName  string
	fromEmail string
	sendMail  smtpSendMailFunc
	wait      retryWaitFunc
}

func NewSMTPMailer(cfg SMTPConfig) (*SMTPMailer, error) {
	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.Username = strings.TrimSpace(cfg.Username)
	cfg.FromName = strings.TrimSpace(cfg.FromName)
	cfg.FromEmail = strings.TrimSpace(cfg.FromEmail)

	if cfg.Host == "" {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("SMTP_PORT must be between 1 and 65535")
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("SMTP_USERNAME is required")
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD is required")
	}
	if err := validateDisplayName(cfg.FromName); err != nil {
		return nil, fmt.Errorf("invalid MAIL_FROM_NAME: %w", err)
	}
	if err := validateAddress(cfg.FromEmail); err != nil {
		return nil, fmt.Errorf("invalid MAIL_FROM_EMAIL: %w", err)
	}

	return &SMTPMailer{
		host:      cfg.Host,
		addr:      cfg.Host + ":" + strconv.Itoa(cfg.Port),
		username:  cfg.Username,
		password:  cfg.Password,
		fromName:  cfg.FromName,
		fromEmail: cfg.FromEmail,
		sendMail:  smtp.SendMail,
		wait:      waitForRetry,
	}, nil
}

func (s *SMTPMailer) Send(ctx context.Context, templateFile, username, email string, data any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.ContainsAny(username, "\r\n") {
		return fmt.Errorf("recipient name contains invalid newline characters")
	}
	if err := validateAddress(email); err != nil {
		return fmt.Errorf("invalid recipient email: %w", err)
	}

	rendered, err := renderTemplate(templateFile, data)
	if err != nil {
		return err
	}
	message := s.buildMessage(username, email, rendered)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	recipients := []string{email}

	var retryErr error
	for attempt := range maxRetries {
		if err := ctx.Err(); err != nil {
			return err
		}

		retryErr = s.sendMail(s.addr, auth, s.fromEmail, recipients, message)
		if retryErr == nil {
			return nil
		}
		if attempt < maxRetries-1 {
			if err := s.wait(ctx, time.Second*time.Duration(attempt+1)); err != nil {
				return err
			}
		}
	}

	return fmt.Errorf("failed to send email to %s after %d attempts: %w", email, maxRetries, retryErr)
}

func (s *SMTPMailer) buildMessage(username, email string, rendered renderedMessage) []byte {
	from := (&mail.Address{Name: s.fromName, Address: s.fromEmail}).String()
	to := (&mail.Address{Name: username, Address: email}).String()

	var message bytes.Buffer
	fmt.Fprintf(&message, "From: %s\r\n", from)
	fmt.Fprintf(&message, "To: %s\r\n", to)
	fmt.Fprintf(&message, "Subject: %s\r\n", rendered.subject)
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	message.WriteString("\r\n")
	message.WriteString(rendered.body)

	return message.Bytes()
}

func validateAddress(address string) error {
	if strings.ContainsAny(address, "\r\n") {
		return fmt.Errorf("address contains invalid newline characters")
	}
	parsed, err := mail.ParseAddress(address)
	if err != nil {
		return err
	}
	if parsed.Address != address {
		return fmt.Errorf("address must not include a display name")
	}
	return nil
}

func waitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
