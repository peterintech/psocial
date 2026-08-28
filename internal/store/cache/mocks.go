package cache

import (
	"context"

	"github.com/peterintech/psocial/internal/store"
	"github.com/stretchr/testify/mock"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	mock.Mock
}

func (s *MockUserStore) GetByID(ctx context.Context, userID string) (*store.User, error) {
	args := s.Called(userID)
	return nil, args.Error(1)
}

func (s *MockUserStore) Set(ctx context.Context, user *store.User) error {
	args := s.Called(user)
	return args.Error(0)
}
