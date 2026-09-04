package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	UserID    string    `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments,omitempty"`
	User      *User     `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, tags, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.UserID).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	return err
}

func (s *PostStore) GetByID(ctx context.Context, postID string) (*Post, error) {
	query := `SELECT p.id, p.content, p.title, p.tags, p.version, p.user_id,
		p.created_at, p.updated_at, u.username
		FROM posts p JOIN users u ON u.id = p.user_id WHERE p.id = $1`

	post := &Post{User: &User{}}
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&post.ID, &post.Content, &post.Title, pq.Array(&post.Tags), &post.Version, &post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.User.Username)

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
	query := `UPDATE posts SET content = $1, title = $2, tags = $3, updated_at = CURRENT_TIMESTAMP, version = version + 1 WHERE id = $4 AND version = $5 RETURNING updated_at, version`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.ID, post.Version).Scan(&post.UpdatedAt, &post.Version)

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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s *PostStore) GetUserFeeds(ctx context.Context, userID string, fq PaginatedFeedQuery) ([]*PostWithMetadata, error) {
	query := `
		SELECT
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN followers f ON f.user_id = p.user_id AND f.follower_id = $1
		WHERE
			(f.follower_id = $1 OR p.user_id = $1) AND
			(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			(p.tags @> $5 OR $5 = '{}')
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := make([]*PostWithMetadata, 0)

	for rows.Next() {
		post := &PostWithMetadata{}
		post.User = &User{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Version, pq.Array(&post.Tags), &post.User.Username, &post.CommentCount)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}
