package cache

import (
	"context"

	"github.com/peterintech/psocial/internal/store"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	users map[string]*store.User
}

func (s *MockUserStore) GetByID(ctx context.Context, userID string) (*store.User, error) {
	if s.users == nil {
		s.users = make(map[string]*store.User)
	}
	user, exists := s.users[userID]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (s *MockUserStore) Set(ctx context.Context, user *store.User) error {
	if s.users == nil {
		s.users = make(map[string]*store.User)
	}
	s.users[user.ID] = user
	return nil
}
