package mailer

import (
	"context"
	"embed"
	"fmt"
	"strings"
)

const (
	FromName            = "Psocial"
	FromEmail           = "peter@cloverkrafts.com"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(ctx context.Context, templateFile, username, email string, data any) error
}

func validateDisplayName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("display name cannot be empty")
	}
	if strings.ContainsAny(name, "\r\n") {
		return fmt.Errorf("display name contains invalid newline characters")
	}
	return nil
}
