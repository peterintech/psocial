package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/peterintech/psocial/internal/store"
)

type Storage struct {
	Users interface {
		GetByID(ctx context.Context, userID string) (*store.User, error)
		Set(ctx context.Context, user *store.User) error
		// Delete(ctx context.Context, userID string) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb},
	}
}
