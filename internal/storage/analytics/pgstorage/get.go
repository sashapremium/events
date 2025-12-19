package pgstorage

import (
	"context"
	"fmt"

	analyticsmodels "github.com/sashapremium/events/internal/pb/models"
	svc "github.com/sashapremium/events/internal/services/analyticsService"
)

func (s *PGStorage) GetPostTotals(ctx context.Context, postID uint64) (svc.PostTotals, error) {
	var viewsTotal, viewsUnique, likes, comments, reposts int64

	err := s.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE type = 'view')                    AS views_total,
			COUNT(DISTINCT user_hash) FILTER (WHERE type = 'view')   AS views_unique,
			COUNT(*) FILTER (WHERE type = 'like')                    AS likes,
			COUNT(*) FILTER (WHERE type = 'comment')                 AS comments,
			COUNT(*) FILTER (WHERE type = 'repost')                  AS reposts
		FROM content_events
		WHERE content_id = $1
	`, postID).Scan(&viewsTotal, &viewsUnique, &likes, &comments, &reposts)
	if err != nil {
		return svc.PostTotals{}, fmt.Errorf("select post totals: %w", err)
	}

	return svc.PostTotals{
		Views:       uint64(viewsTotal),
		Likes:       uint64(likes),
		Comments:    uint64(comments),
		Reposts:     uint64(reposts),
		UniqueUsers: uint64(viewsUnique),
	}, nil
}

func (s *PGStorage) GetTopPostsByType(
	ctx context.Context,
	eventType string,
	limit uint64,
) ([]svc.TopItem, error) {
	rows, err := s.db.Query(ctx, `
		SELECT content_id, COUNT(*) AS cnt
		FROM content_events
		WHERE type = $1
		GROUP BY content_id
		ORDER BY cnt DESC
		LIMIT $2
	`, eventType, limit)
	if err != nil {
		return nil, fmt.Errorf("select top posts: %w", err)
	}
	defer rows.Close()

	out := make([]svc.TopItem, 0, limit)
	for rows.Next() {
		var contentID int64
		var cnt int64
		if err := rows.Scan(&contentID, &cnt); err != nil {
			return nil, fmt.Errorf("scan top posts: %w", err)
		}
		if contentID < 0 {
			continue
		}
		out = append(out, svc.TopItem{
			PostID: uint64(contentID),
			Value:  uint64(cnt),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows top posts err: %w", err)
	}

	return out, nil
}

func (s *PGStorage) GetAuthorStats(ctx context.Context, authorID uint64) (*analyticsmodels.AuthorStatsModel, error) {
	var posts int64
	var totalViews, totalLikes, totalComments, totalReposts int64

	err := s.db.QueryRow(ctx, `
		SELECT
			COUNT(DISTINCT content_id) FILTER (WHERE author_id = $1)        AS posts,
			COUNT(*) FILTER (WHERE author_id = $1 AND type = 'view')        AS total_views,
			COUNT(*) FILTER (WHERE author_id = $1 AND type = 'like')        AS total_likes,
			COUNT(*) FILTER (WHERE author_id = $1 AND type = 'comment')     AS total_comments,
			COUNT(*) FILTER (WHERE author_id = $1 AND type = 'repost')      AS total_reposts
		FROM content_events
	`, authorID).Scan(&posts, &totalViews, &totalLikes, &totalComments, &totalReposts)
	if err != nil {
		return nil, fmt.Errorf("select author stats: %w", err)
	}

	return &analyticsmodels.AuthorStatsModel{
		AuthorId:      fmt.Sprintf("%d", authorID),
		Posts:         uint64(posts),
		TotalViews:    totalViews,
		TotalLikes:    totalLikes,
		TotalComments: totalComments,
		TotalReposts:  totalReposts,
	}, nil
}
