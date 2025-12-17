package pgstorage

import (
	"context"
	"fmt"

	svc "github.com/sashapremium/events/internal/services/analyticsService"
)

func (s *PGStorage) GetPostTotals(ctx context.Context, postID uint64) (svc.PostTotals, error) {
	var viewsTotal, viewsUnique, likes, comments, reposts int64

	err := s.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE type = 'view')                                   AS views_total,
			COUNT(DISTINCT user_hash) FILTER (WHERE type = 'view')                  AS views_unique,
			COUNT(*) FILTER (WHERE type = 'like')                                   AS likes,
			COUNT(*) FILTER (WHERE type = 'comment')                                AS comments,
			COUNT(*) FILTER (WHERE type = 'repost')                                 AS reposts
		FROM content_events
		WHERE content_id = $1
	`, int64(postID)).Scan(&viewsTotal, &viewsUnique, &likes, &comments, &reposts)

	if err != nil {
		return svc.PostTotals{}, fmt.Errorf("select post totals: %w", err)
	}

	return svc.PostTotals{
		Views:       viewsTotal,
		Likes:       likes,
		Comments:    comments,
		Reposts:     reposts,
		UniqueUsers: viewsUnique,
	}, nil
}
