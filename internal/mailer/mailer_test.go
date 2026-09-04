package mailer

import (
	"strings"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	data := struct {
		Username      string
		ActivationURL string
	}{
		Username:      `<script>alert("x")</script>`,
		ActivationURL: "https://example.com/activate?token=test",
	}

	rendered, err := renderTemplate(UserWelcomeTemplate, data)
	if err != nil {
		t.Fatalf("renderTemplate returned an error: %v", err)
	}
	if rendered.subject != "Finish Registration with Psocial" {
		t.Fatalf("unexpected subject %q", rendered.subject)
	}
	if strings.Contains(rendered.body, "<script>") {
		t.Fatal("expected HTML template data to be escaped")
	}
	if !strings.Contains(rendered.body, "&lt;script&gt;") {
		t.Fatal("escaped username was not present in rendered body")
	}
}

func TestRenderTemplateRejectsMissingTemplate(t *testing.T) {
	if _, err := renderTemplate("missing.tmpl", nil); err == nil {
		t.Fatal("expected a missing template error")
	}
}
