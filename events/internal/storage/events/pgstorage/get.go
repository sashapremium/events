package pgstorage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	eventmodel "github.com/sashapremium/events/internal/models"
)

var ErrPostNotFound = errors.New("post not found")

func (s *PGStorage) GetPost(ctx context.Context, id uint64) (*eventmodel.PostInfo, error) {
	const q = `
		SELECT id, title, author_id, category, content, published_at
		FROM posts
		WHERE id = $1
	`
	var p eventmodel.PostInfo

	err := s.db.QueryRow(ctx, q, id).Scan(
		&p.ID,
		&p.Title,
		&p.AuthorID,
		&p.Category,
		&p.Content,
		&p.PublishedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &p, nil
}
func (s *PGStorage) GetPostAuthorID(ctx context.Context, postID uint64) (uint64, error) {
	const q = `SELECT author_id FROM posts WHERE id = $1`
	var authorID uint64
	err := s.db.QueryRow(ctx, q, postID).Scan(&authorID)
	return authorID, err
}
