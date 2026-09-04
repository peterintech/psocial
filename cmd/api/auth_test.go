package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/peterintech/psocial/internal/store"
	"go.uber.org/zap"
)

type registrationMailer struct {
	err   error
	calls int
}

func (m *registrationMailer) Send(context.Context, string, string, string, any) error {
	m.calls++
	return m.err
}

type registrationUserStore struct {
	created bool
	deleted bool
}

func (s *registrationUserStore) Create(context.Context, *sql.Tx, *store.User) error {
	return nil
}

func (s *registrationUserStore) GetByID(context.Context, string) (*store.User, error) {
	return nil, store.ErrNotFound
}

func (s *registrationUserStore) GetByEmail(context.Context, string) (*store.User, error) {
	return nil, store.ErrNotFound
}

func (s *registrationUserStore) CreateAndInvite(_ context.Context, user *store.User, _ string, _ time.Duration) error {
	s.created = true
	user.ID = "test-user-id"
	return nil
}

func (s *registrationUserStore) Activate(context.Context, string) error {
	return nil
}

func (s *registrationUserStore) Delete(context.Context, string) error {
	s.deleted = true
	return nil
}

func TestRegisterUserMailDelivery(t *testing.T) {
	tests := []struct {
		name        string
		mailError   error
		wantStatus  int
		wantDeleted bool
	}{
		{name: "successful delivery", wantStatus: http.StatusCreated},
		{name: "failed delivery rolls back user", mailError: errors.New("delivery failed"), wantStatus: http.StatusInternalServerError, wantDeleted: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &registrationUserStore{}
			mailClient := &registrationMailer{err: tt.mailError}
			app := &application{
				config: config{
					frontendURL: "http://localhost:5173/",
					mail: mailConfig{
						provider: "smtp",
						exp:      72 * time.Hour,
					},
				},
				store:  store.Storage{Users: users},
				logger: zap.NewNop().Sugar(),
				mailer: mailClient,
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(`{
				"username":"peter",
				"email":"peter@example.com",
				"password":"password123"
			}`))
			recorder := httptest.NewRecorder()

			app.registerUserHandler(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, recorder.Code)
			}
			if !users.created {
				t.Fatal("expected the user to be created")
			}
			if mailClient.calls != 1 {
				t.Fatalf("expected one mail delivery attempt, got %d", mailClient.calls)
			}
			if users.deleted != tt.wantDeleted {
				t.Fatalf("expected deleted=%t, got %t", tt.wantDeleted, users.deleted)
			}
		})
	}
}
