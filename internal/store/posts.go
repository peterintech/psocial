package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	UserID    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments,omitempty"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, tags, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.UserID).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	return err
}

func (s *PostStore) GetByID(ctx context.Context, postID string) (*Post, error) {
	query := `SELECT id, content, title, tags, user_id, created_at, updated_at FROM posts WHERE id = $1`

	post := &Post{}
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&post.ID, &post.Content, &post.Title, pq.Array(&post.Tags), &post.UserID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts SET content = $1, title = $2, tags = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4 RETURNING updated_at`

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.ID).Scan(&post.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) Delete(ctx context.Context, postID string) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
