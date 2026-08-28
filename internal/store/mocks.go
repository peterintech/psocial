package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	// You can add fields here to store mock data or track method calls if needed
}

func (m *MockUserStore) GetByID(ctx context.Context, userID string) (*User, error) {
	// Implement mock behavior for GetByID
	return nil, nil
}

func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	// Implement mock behavior for Activate
	return nil
}
func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	// Implement mock behavior for CreateAndInvite
	return nil
}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	// Implement mock behavior for Create
	return nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	// Implement mock behavior for GetByEmail
	return nil, nil
}

func (m *MockUserStore) Delete(ctx context.Context, userID string) error {
	// Implement mock behavior for Delete
	return nil
}
