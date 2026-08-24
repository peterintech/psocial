package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (content, post_id, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, ctxDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, comment.Content, comment.PostID, comment.UserID).Scan(&comment.ID, &comment.CreatedAt)

	return err
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID string) ([]Comment, error) {
	query := `SELECT c.id, c.content, c.post_id, c.user_id, c.created_at, u.username
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.PostID, &c.UserID, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
