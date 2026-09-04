package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type renderedMessage struct {
	subject string
	body    string
}

func renderTemplate(templateFile string, data any) (renderedMessage, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return renderedMessage{}, fmt.Errorf("failed to parse mail template: %w", err)
	}

	var subject, body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&subject, "subject", data); err != nil {
		return renderedMessage{}, fmt.Errorf("failed to render mail subject: %w", err)
	}
	if err := tmpl.ExecuteTemplate(&body, "body", data); err != nil {
		return renderedMessage{}, fmt.Errorf("failed to render mail body: %w", err)
	}

	renderedSubject := strings.TrimSpace(subject.String())
	if renderedSubject == "" {
		return renderedMessage{}, fmt.Errorf("mail subject cannot be empty")
	}

	return renderedMessage{subject: renderedSubject, body: body.String()}, nil
}
