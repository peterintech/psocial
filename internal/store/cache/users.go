package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/peterintech/psocial/internal/store"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStore) GetByID(ctx context.Context, userID string) (*store.User, error) {
	cacheKey := "user-" + userID

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err = json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := "user-" + user.ID
	jsonRes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx, cacheKey, jsonRes, UserExpTime).Err()
}
