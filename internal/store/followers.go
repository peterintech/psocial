package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     string `json:"user_id"`
	FollowerID string `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID string) error {
	query := `INSERT INTO followers (follower_id, user_id)
		VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followerID, userID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}

	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, unfollowedID, userID string) error {
	query := `DELETE FROM followers WHERE follower_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, unfollowedID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user %s: %w, not followed", userID, ErrNotFound)
	}

	return nil
}

func (s *FollowerStore) IsFollowing(ctx context.Context, followerID, userID string) (bool, error) {
	const query = `SELECT EXISTS (
		SELECT 1 FROM followers WHERE follower_id = $1 AND user_id = $2
	)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var following bool
	if err := s.db.QueryRowContext(ctx, query, followerID, userID).Scan(&following); err != nil {
		return false, err
	}
	return following, nil
}
