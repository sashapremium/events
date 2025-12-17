package pgstorage

import (
	"context"
	"fmt"
	"time"

	analyticsmodels "github.com/sashapremium/events/internal/pb/models"
)

func (s *PGStorage) GetTopPostsByType(ctx context.Context, eventType string, from, to time.Time, limit uint64) ([]TopPostItem, error) {
	rows, err := s.db.Query(ctx, `
		SELECT content_id, COUNT(*) AS cnt
		FROM content_events
		WHERE type = $1 AND at >= $2 AND at < $3
		GROUP BY content_id
		ORDER BY cnt DESC
		LIMIT $4
	`, eventType, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("select top posts: %w", err)
	}
	defer rows.Close()

	out := make([]TopPostItem, 0, limit)
	for rows.Next() {
		var contentID string
		var cnt int64
		if err := rows.Scan(&contentID, &cnt); err != nil {
			return nil, fmt.Errorf("scan top posts: %w", err)
		}
		out = append(out, TopPostItem{ContentID: contentID, Count: uint64(cnt)})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows top posts err: %w", err)
	}

	return out, nil
}

func (s *PGStorage) GetTopAuthorsByType(ctx context.Context, eventType string, from, to time.Time, limit uint64) ([]TopAuthorItem, error) {
	rows, err := s.db.Query(ctx, `
		SELECT author_id, COUNT(*) AS cnt
		FROM content_events
		WHERE author_id IS NOT NULL AND type = $1 AND at >= $2 AND at < $3
		GROUP BY author_id
		ORDER BY cnt DESC
		LIMIT $4
	`, eventType, from, to, limit)
	if err != nil {
		return nil, fmt.Errorf("select top authors: %w", err)
	}
	defer rows.Close()

	out := make([]TopAuthorItem, 0, limit)
	for rows.Next() {
		var authorID string
		var cnt int64
		if err := rows.Scan(&authorID, &cnt); err != nil {
			return nil, fmt.Errorf("scan top authors: %w", err)
		}
		out = append(out, TopAuthorItem{AuthorID: authorID, Count: uint64(cnt)})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows top authors err: %w", err)
	}

	return out, nil
}

func (s *PGStorage) GetPostStats(ctx context.Context, contentID string, from, to time.Time) (PostStats, error) {
	st := PostStats{
		ContentID:    contentID,
		CalculatedAt: time.Now().UTC(),
	}

	err := s.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE type = 'view')                                   AS views_total,
			COUNT(DISTINCT user_hash) FILTER (WHERE type = 'view')                  AS views_unique,
			COUNT(*) FILTER (WHERE type = 'like')                                   AS likes,
			COUNT(*) FILTER (WHERE type = 'repost')                                 AS reposts,
			COUNT(*) FILTER (WHERE type = 'comment')                                AS comments
		FROM content_events
		WHERE content_id = $1 AND at >= $2 AND at < $3
	`, contentID, from, to).Scan(&st.ViewsTotal, &st.ViewsUnique, &st.Likes, &st.Reposts, &st.Comments)
	if err != nil {
		return PostStats{}, fmt.Errorf("select post stats: %w", err)
	}

	rows, err := s.db.Query(ctx, `
		SELECT date_trunc('day', at) AS day, COUNT(*) AS cnt
		FROM content_events
		WHERE content_id = $1 AND type = 'view' AND at >= $2 AND at < $3
		GROUP BY day
		ORDER BY day ASC
	`, contentID, from, to)
	if err != nil {
		return PostStats{}, fmt.Errorf("select daily views: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var day time.Time
		var cnt int64
		if err := rows.Scan(&day, &cnt); err != nil {
			return PostStats{}, fmt.Errorf("scan daily views: %w", err)
		}
		st.DailyViews = append(st.DailyViews, DailyCount{Day: day, Count: uint64(cnt)})
	}
	if err := rows.Err(); err != nil {
		return PostStats{}, fmt.Errorf("rows daily views err: %w", err)
	}

	return st, nil
}

func (s *PGStorage) GetAuthorStats(ctx context.Context, authorID string) (*analyticsmodels.AuthorStatsModel, error) {
	var posts int64
	var totalViews, totalLikes, totalComments, totalReposts int64

	err := s.db.QueryRow(ctx, `
		SELECT
			COUNT(DISTINCT content_id) FILTER (WHERE author_id = $1)                                                AS posts,
			COUNT(*) FILTER (WHERE author_id = $1 AND type = 'view')                                                AS total_views,
			COUNT(*) FILTER (WHERE author_id = $1 AND type = 'like')                                                AS total_likes,
			COUNT(*) FILTER (WHERE author_id = $1 AND type = 'comment')                                             AS total_comments,
			COUNT(*) FILTER (WHERE author_id = $1 AND type = 'repost')                                              AS total_reposts
		FROM content_events
	`, authorID).Scan(&posts, &totalViews, &totalLikes, &totalComments, &totalReposts)

	if err != nil {
		return nil, fmt.Errorf("select author stats: %w", err)
	}

	out := &analyticsmodels.AuthorStatsModel{
		AuthorId:      authorID,
		Posts:         uint64(posts),
		TotalViews:    totalViews,
		TotalLikes:    totalLikes,
		TotalComments: totalComments,
		TotalReposts:  totalReposts,
	}

	return out, nil
}
