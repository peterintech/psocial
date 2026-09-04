package mailer

import (
	"context"
	"embed"
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
