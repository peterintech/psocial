package store

import (
	"context"
	"database/sql"
)

type UserStore struct {
	db *sql.DB
}

func (us *UserStore) Create(ctx context.Context) error {
	// Implement the logic to create a user in the database
	return nil
}
