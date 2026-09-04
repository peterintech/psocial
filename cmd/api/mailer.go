package main

import (
	"fmt"
	"strings"

	"github.com/peterintech/psocial/internal/mailer"
)

func resolveNewMailer(cfg mailConfig) (mailer.Client, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.provider)) {
	case "sendgrid":
		return mailer.NewSendGridMailer(cfg.fromName, cfg.fromEmail, cfg.sendGrid.apiKey)
	case "smtp":
		return mailer.NewSMTPMailer(mailer.SMTPConfig{
			Host:      cfg.smtp.host,
			Port:      cfg.smtp.port,
			Username:  cfg.smtp.username,
			Password:  cfg.smtp.password,
			FromName:  cfg.fromName,
			FromEmail: cfg.fromEmail,
		})
	default:
		return nil, fmt.Errorf("unsupported MAIL_PROVIDER %q: expected sendgrid or smtp", cfg.provider)
	}
}
