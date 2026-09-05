package mailer

import "context"

type mockMailer struct{}

func (m mockMailer) Send(ctx context.Context, templateFile, username, email string, data any) error {
	return nil
}

func NewMockMailer() Client {
	return &mockMailer{}
}
