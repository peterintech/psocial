package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrConflict          = errors.New("record already exists")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, string) (*Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, string) error
		GetUserFeeds(ctx context.Context, userID string, fq PaginatedFeedQuery) ([]*PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, string) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, string) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID string) error
		Unfollow(ctx context.Context, unfollowedID, userID string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
